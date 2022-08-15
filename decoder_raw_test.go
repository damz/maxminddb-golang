package maxminddb

import (
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRawDecoder(t *testing.T) {
	require := require.New(t)

	reader, err := Open(testFile("MaxMind-DB-test-decoder.mmdb"))
	require.NoError(err)

	offset, err := reader.LookupOffset(net.ParseIP("::1.1.1.0"))
	require.NoError(err)

	d := reader.Decoder(offset)

	err = d.DecodeMap(func(key string, val *Decoder) error {
		switch key {
		case "array":
			var values []uint32
			err := val.DecodeSlice(func(val *Decoder) error {
				vv, err := val.DecodeUInt32()
				require.NoError(err)

				values = append(values, vv)
				return nil
			})
			require.NoError(err)

			require.Equal([]uint32{1, 2, 3}, values)

		case "boolean":
			vv, err := val.DecodeBool()
			require.NoError(err)
			require.Equal(true, vv)

		case "bytes":
			vv, err := val.DecodeBytes()
			require.NoError(err)
			require.Equal([]byte{0x00, 0x00, 0x00, 0x2a}, vv)

		case "double":
			vv, err := val.DecodeFloat64()
			require.NoError(err)
			require.Equal(float64(42.123456), vv)

		case "float":
			vv, err := val.DecodeFloat32()
			require.NoError(err)
			require.Equal(float32(1.1), vv)

		case "int32":
			vv, err := val.DecodeInt32()
			require.NoError(err)
			require.Equal(int32(-268435456), vv)

		case "map":
			var keys []string
			err := val.DecodeMap(func(key string, val *Decoder) error {
				keys = append(keys, key)

				if key == "mapX" {
					var subKeys []string
					err := val.DecodeMap(func(key string, val *Decoder) error {
						subKeys = append(subKeys, key)

						switch key {
						case "arrayX":
							var values []uint32
							err := val.DecodeSlice(func(val *Decoder) error {
								vv, err := val.DecodeUInt32()
								require.NoError(err)

								values = append(values, vv)
								return nil
							})
							require.NoError(err)

							require.Equal([]uint32{7, 8, 9}, values)

						case "utf8_stringX":
							vv, err := val.DecodeString()
							require.NoError(err)
							require.Equal("hello", vv)

						default:
							return fmt.Errorf("unexpected key: %#v", key)
						}

						return nil
					})
					require.NoError(err)
					require.Equal([]string{"arrayX", "utf8_stringX"}, subKeys)
				}
				return nil
			})
			require.NoError(err)
			require.Equal([]string{"mapX"}, keys)

		case "uint16":
			vv, err := val.DecodeUInt16()
			require.NoError(err)
			require.Equal(uint16(100), vv)

		case "uint32":
			vv, err := val.DecodeUInt32()
			require.NoError(err)
			require.Equal(uint32(268435456), vv)

		case "uint64":
			vv, err := val.DecodeUInt64()
			require.NoError(err)
			require.Equal(uint64(1152921504606846976), vv)

		case "uint128":
			hi, lo, err := val.DecodeUInt128()
			require.NoError(err)
			require.Equal(uint64(0x100000000000000), hi)
			require.Equal(uint64(0x000000000000000), lo)

		case "utf8_string":
			vv, err := val.DecodeString()
			require.NoError(err)
			require.Equal("unicode! ☯ - ♫", vv)

		default:
			log.Printf("Key: %#v", key)
		}
		return nil
	})
	require.NoError(err)
}
