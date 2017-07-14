package encodings

import (
	"io"
	"vncproxy/common"
	"vncproxy/logger"
)

const (
	HextileRaw                 = 1
	HextileBackgroundSpecified = 2
	HextileForegroundSpecified = 4
	HextileAnySubrects         = 8
	HextileSubrectsColoured    = 16
)

type HextileEncoding struct {
	//Colors []Color
	bytes []byte
}

func (z *HextileEncoding) Type() int32 {
	return 5
}
func (z *HextileEncoding) WriteTo(w io.Writer) (n int, err error) {
	return w.Write(z.bytes)
}
func (z *HextileEncoding) Read(pixelFmt *common.PixelFormat, rect *common.Rectangle, r *common.RfbReadHelper) (common.Encoding, error) {
	bytesPerPixel := int(pixelFmt.BPP) / 8

	r.StartByteCollection()
	defer func() {
		z.bytes = r.EndByteCollection()
	}()

	for ty := rect.Y; ty < rect.Y+rect.Height; ty += 16 {
		th := 16
		if rect.Y+rect.Height-ty < 16 {
			th = int(rect.Y) + int(rect.Height) - int(ty)
		}

		for tx := rect.X; tx < rect.X+rect.Width; tx += 16 {
			tw := 16
			if rect.X+rect.Width-tx < 16 {
				tw = int(rect.X) + int(rect.Width) - int(tx)
			}

			//handle Hextile Subrect(tx, ty, tw, th):
			subencoding, err := r.ReadUint8()
			//logger.Debugf("hextile reader tile: (%d,%d) subenc=%d\n", ty, tx, subencoding)
			if err != nil {
				logger.Errorf("HextileEncoding.Read: error in hextile reader: %v", err)
				return nil, err
			}

			if (subencoding & HextileRaw) != 0 {
				//ReadRawRect(c, rect, r)
				r.ReadBytes(tw * th * bytesPerPixel)
				//logger.Debug("hextile reader: HextileRaw\n")
				continue
			}
			if (subencoding & HextileBackgroundSpecified) != 0 {
				r.ReadBytes(int(bytesPerPixel))
			}
			if (subencoding & HextileForegroundSpecified) != 0 {
				r.ReadBytes(int(bytesPerPixel))
			}
			if (subencoding & HextileAnySubrects) == 0 {
				//logger.Debug("hextile reader: no Subrects")
				continue
			}

			nSubrects, err := r.ReadUint8()
			if err != nil {
				return nil, err
			}
			bufsize := int(nSubrects) * 2
			if (subencoding & HextileSubrectsColoured) != 0 {
				bufsize += int(nSubrects) * int(bytesPerPixel)
			}
			//byte[] buf = new byte[bufsize];
			r.ReadBytes(bufsize)
		}
	}

	// len, _ := readUint32(c.c)
	// _, err := readBytes(c.c, int(len))

	// if err != nil {
	// 	return nil, err
	// }
	return z, nil
}
