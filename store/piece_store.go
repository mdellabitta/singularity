package store

import (
	"context"
	"io"

	"github.com/data-preservation-programs/go-singularity/datasource"
	"github.com/data-preservation-programs/go-singularity/model"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-varint"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

type PieceBlock interface {
	GetPieceOffset() uint64
}

type ItemBlockMetadata struct {
	PieceOffset uint64 `json:"PieceOffset"`
	Varint      []byte `json:"Varint"`
	Cid         []byte `json:"Cid"`
	ItemOffset  uint64 `json:"itemOffset"`
	ItemLength  uint64 `json:"itemLength"`
}

func (i ItemBlockMetadata) GetPieceOffset() uint64 {
	return i.PieceOffset
}
func (i ItemBlockMetadata) CidOffset() uint64 {
	return i.PieceOffset + uint64(len(i.Varint))
}
func (i ItemBlockMetadata) BlockOffset() uint64 {
	return i.PieceOffset + uint64(len(i.Varint)) + uint64(len(i.Cid))
}
func (i ItemBlockMetadata) EndOffset() uint64 {
	return i.PieceOffset + uint64(len(i.Varint)) + uint64(len(i.Cid)) + uint64(i.ItemLength)
}
func (i ItemBlockMetadata) Length() int {
	return len(i.Varint) + len(i.Cid) + int(i.ItemLength)
}

type RawBlock struct {
	PieceOffset uint64 `json:"pieceOffset"`
	Varint      []byte `json:"varint"`
	Cid         []byte `json:"cid"`
	BlockData   []byte `json:"blockData"`
}

func (r RawBlock) GetPieceOffset() uint64 {
	return r.PieceOffset
}
func (r RawBlock) CidOffset() uint64 {
	return r.PieceOffset + uint64(len(r.Varint))
}
func (r RawBlock) BlockOffset() uint64 {
	return r.PieceOffset + uint64(len(r.Varint)) + uint64(len(r.Cid))
}
func (r RawBlock) EndOffset() uint64 {
	return r.PieceOffset + uint64(len(r.Varint)) + uint64(len(r.Cid)) + uint64(len(r.BlockData))
}

func (r RawBlock) Length() int {
	return len(r.Varint) + len(r.Cid) + len(r.BlockData)
}

type ItemBlock struct {
	PieceOffset   uint64              `json:"pieceOffset"`
	SourceHandler datasource.Handler  `json:"-"`
	Item          *model.Item         `json:"item"`
	Meta          []ItemBlockMetadata `json:"meta"`
}

func (i ItemBlock) GetPieceOffset() uint64 {
	return i.PieceOffset
}

type PieceReader struct {
	ctx          context.Context
	Blocks       []PieceBlock `json:"blocks"`
	reader       io.ReadCloser
	pos          uint64
	blockID      int
	innerBlockID int
	blockOffset  uint64
	Header       []byte `json:"header"`
}

func (pr *PieceReader) MakeCopy(ctx context.Context, offset uint64) (*PieceReader, error) {
	newReader := &PieceReader{
		ctx:    ctx,
		Blocks: pr.Blocks,
		reader: nil,
		pos:    offset,
		Header: pr.Header,
	}

	if offset < uint64(len(pr.Header)) {
		return newReader, nil
	}

	index, _ := slices.BinarySearchFunc(
		pr.Blocks, offset, func(b PieceBlock, o uint64) int {
			return int(b.GetPieceOffset() - o)
		},
	)
	newReader.blockID = index
	switch block := pr.Blocks[index].(type) {
	case RawBlock:
		newReader.blockOffset = offset - block.GetPieceOffset()
	case ItemBlock:
		innerIndex, _ := slices.BinarySearchFunc(
			block.Meta, offset, func(b ItemBlockMetadata, o uint64) int {
				return int(b.GetPieceOffset() - o)
			},
		)
		newReader.innerBlockID = innerIndex
		newReader.blockOffset = offset - block.Meta[innerIndex].GetPieceOffset()
	}

	return newReader, nil
}

func NewPieceReader(
	ctx context.Context,
	car model.Car,
	carBlocks []model.CarBlock,
	resolver datasource.HandlerResolver,
) (
	*PieceReader,
	error,
) {
	// Sanitize carBlocks
	if len(carBlocks) == 0 {
		return nil, errors.New("no Blocks provided")
	}

	if carBlocks[0].CarOffset != uint64(len(car.Header)) {
		return nil, errors.New("first block must start at car Header")
	}

	lastBlock := carBlocks[len(carBlocks)-1]
	if lastBlock.CarOffset+lastBlock.CarBlockLength != car.FileSize {
		return nil, errors.New("last block must end at car footer")
	}

	for i := 0; i < len(carBlocks)-1; i++ {
		if carBlocks[i].CarOffset+carBlocks[i].CarBlockLength != carBlocks[i+1].CarOffset {
			return nil, errors.New("Blocks must be contiguous")
		}
		if carBlocks[i].RawBlock == nil && (carBlocks[i].Item == nil || carBlocks[i].Source == nil) {
			return nil, errors.New("block must be either raw or Item, and the Item/source needs to be preloaded")
		}
	}

	// Combine nearby clocks with same Item
	blocks := make([]PieceBlock, 0)
	var lastItemBlock *ItemBlock
	for _, carBlock := range carBlocks {
		if lastItemBlock != nil && (carBlock.RawBlock != nil || lastItemBlock.Item.ID != carBlock.Item.ID) {
			blocks = append(blocks, *lastItemBlock)
			lastItemBlock = nil
		}
		if carBlock.RawBlock != nil {
			blocks = append(
				blocks, RawBlock{
					PieceOffset: carBlock.CarOffset,
					Varint:      varint.ToUvarint(carBlock.Varint),
					Cid:         cid.MustParse(carBlock.CID).Bytes(),
					BlockData:   carBlock.RawBlock,
				},
			)
			continue
		}
		if lastItemBlock == nil {
			handler, err := resolver.GetHandler(*carBlock.Source)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get handler")
			}
			lastItemBlock = &ItemBlock{
				PieceOffset:   carBlock.CarOffset,
				SourceHandler: handler,
				Item:          carBlock.Item,
				Meta: []ItemBlockMetadata{
					{
						PieceOffset: carBlock.CarOffset,
						Varint:      varint.ToUvarint(carBlock.Varint),
						Cid:         cid.MustParse(carBlock.CID).Bytes(),
						ItemOffset:  carBlock.ItemOffset,
						ItemLength:  carBlock.BlockLength,
					},
				},
			}
			continue
		}
		// merge last Item with the new Item
		lastItemBlock.Meta = append(
			lastItemBlock.Meta, ItemBlockMetadata{
				PieceOffset: carBlock.CarOffset,
				Varint:      varint.ToUvarint(carBlock.Varint),
				Cid:         cid.MustParse(carBlock.CID).Bytes(),
				ItemOffset:  carBlock.ItemOffset,
				ItemLength:  carBlock.BlockLength,
			},
		)
	}
	if lastItemBlock != nil {
		blocks = append(blocks, *lastItemBlock)
	}

	return &PieceReader{
		ctx:          ctx,
		Blocks:       blocks,
		reader:       nil,
		pos:          0,
		blockID:      0,
		innerBlockID: 0,
		Header:       car.Header,
	}, nil
}

func (pr *PieceReader) Read(p []byte) (n int, err error) {
	if pr.blockID >= len(pr.Blocks) {
		return 0, io.EOF
	}
	if pr.pos < uint64(len(pr.Header)) {
		copied := copy(p[n:], pr.Header[pr.pos:])
		pr.pos += uint64(copied)
		n += copied
		if n == len(p) {
			return n, nil
		}
	}
	currentBlock := pr.Blocks[pr.blockID]
	if rawBlock, ok := currentBlock.(RawBlock); ok {
		if pr.pos < rawBlock.CidOffset() {
			copied := copy(p[n:], rawBlock.Varint[pr.pos-rawBlock.PieceOffset:])
			pr.pos += uint64(copied)
			n += copied
			if n == len(p) {
				return n, nil
			}
		}
		if pr.pos < rawBlock.BlockOffset() {
			copied := copy(p[n:], rawBlock.Cid[pr.pos-rawBlock.CidOffset():])
			pr.pos += uint64(copied)
			n += copied
			if n == len(p) {
				return n, nil
			}
		}
		if pr.pos < rawBlock.EndOffset() {
			copied := copy(p[n:], rawBlock.BlockData[pr.pos-rawBlock.BlockOffset():])
			pr.pos += uint64(copied)
			n += copied
			if n == len(p) {
				return n, nil
			}
		}
		pr.blockID++
		pr.innerBlockID = 0
		return n, nil
	}

	itemBlock, _ := currentBlock.(ItemBlock)
	innerBlock := itemBlock.Meta[pr.innerBlockID]
	if pr.reader == nil {
		pr.reader, err = itemBlock.SourceHandler.Read(
			pr.ctx,
			itemBlock.Item.Path,
			innerBlock.ItemOffset+pr.blockOffset,
			itemBlock.Item.Size-(innerBlock.ItemOffset+pr.blockOffset),
		)
		if err != nil {
			return 0, errors.Wrap(err, "failed to read Item")
		}
	}
	if pr.pos < innerBlock.CidOffset() {
		copied := copy(p[n:], innerBlock.Varint[pr.pos-innerBlock.PieceOffset:])
		pr.pos += uint64(copied)
		n += copied
		if n == len(p) {
			return n, nil
		}
	}
	if pr.pos < innerBlock.BlockOffset() {
		copied := copy(p[n:], innerBlock.Cid[pr.pos-innerBlock.CidOffset():])
		pr.pos += uint64(copied)
		n += copied
		if n == len(p) {
			return n, nil
		}
	}
	if pr.pos < innerBlock.EndOffset() {
		readTill := min(len(p), n+int(innerBlock.EndOffset()-pr.pos))
		read, err := pr.reader.Read(p[n:readTill])
		n += read
		pr.pos += uint64(read)
		if err != nil && err != io.EOF {
			return n, errors.Wrap(err, "failed to read Item")
		}
		if pr.pos == innerBlock.EndOffset() {
			pr.innerBlockID++
			if pr.innerBlockID >= len(itemBlock.Meta) {
				pr.blockID++
				pr.innerBlockID = 0
				pr.reader.Close()
				pr.reader = nil
			}
		}
		if n == len(p) {
			return n, nil
		}
	}
	return n, nil
}

func min(i int, i2 int) int {
	if i < i2 {
		return i
	}
	return i2
}

func (pr *PieceReader) Close() error {
	if pr.reader == nil {
		return nil
	}
	return pr.reader.Close()
}
