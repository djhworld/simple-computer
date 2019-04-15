package components

import "testing"

func Test2x4Decoder(t *testing.T) {
	decoder := NewDecoder2x4()

	decoder.Update(false, false)
	checkDecoder2x4Output(decoder, []bool{true, false, false, false}, t)
	decoder.Update(false, true)
	checkDecoder2x4Output(decoder, []bool{false, true, false, false}, t)
	decoder.Update(true, false)
	checkDecoder2x4Output(decoder, []bool{false, false, true, false}, t)
	decoder.Update(true, true)
	checkDecoder2x4Output(decoder, []bool{false, false, false, true}, t)
}

func Test3x8Decoder(t *testing.T) {
	decoder := NewDecoder3x8()

	expected := func(onIndex int, size int) []bool {
		result := make([]bool, size)
		result[onIndex] = true
		return result
	}

	decoder.Update(false, false, false)
	checkDecoder3x8Output(decoder, expected(0, 16), t)

	decoder.Update(false, false, true)
	checkDecoder3x8Output(decoder, expected(1, 16), t)

	decoder.Update(false, true, false)
	checkDecoder3x8Output(decoder, expected(2, 16), t)

	decoder.Update(false, true, true)
	checkDecoder3x8Output(decoder, expected(3, 16), t)

	decoder.Update(true, false, false)
	checkDecoder3x8Output(decoder, expected(4, 16), t)

	decoder.Update(true, false, true)
	checkDecoder3x8Output(decoder, expected(5, 16), t)

	decoder.Update(true, true, false)
	checkDecoder3x8Output(decoder, expected(6, 16), t)

	decoder.Update(true, true, true)
	checkDecoder3x8Output(decoder, expected(7, 16), t)
}

func Test4x16Decoder(t *testing.T) {
	decoder := NewDecoder4x16()

	expected := func(onIndex int, size int) []bool {
		result := make([]bool, size)
		result[onIndex] = true
		return result
	}

	decoder.Update(false, false, false, false)
	checkDecoder4x16Output(decoder, expected(0, 16), t)

	decoder.Update(false, false, false, true)
	checkDecoder4x16Output(decoder, expected(1, 16), t)

	decoder.Update(false, false, true, false)
	checkDecoder4x16Output(decoder, expected(2, 16), t)

	decoder.Update(false, false, true, true)
	checkDecoder4x16Output(decoder, expected(3, 16), t)

	decoder.Update(false, true, false, false)
	checkDecoder4x16Output(decoder, expected(4, 16), t)

	decoder.Update(false, true, false, true)
	checkDecoder4x16Output(decoder, expected(5, 16), t)

	decoder.Update(false, true, true, false)
	checkDecoder4x16Output(decoder, expected(6, 16), t)

	decoder.Update(false, true, true, true)
	checkDecoder4x16Output(decoder, expected(7, 16), t)

	decoder.Update(true, false, false, false)
	checkDecoder4x16Output(decoder, expected(8, 16), t)

	decoder.Update(true, false, false, true)
	checkDecoder4x16Output(decoder, expected(9, 16), t)

	decoder.Update(true, false, true, false)
	checkDecoder4x16Output(decoder, expected(10, 16), t)

	decoder.Update(true, false, true, true)
	checkDecoder4x16Output(decoder, expected(11, 16), t)

	decoder.Update(true, true, false, false)
	checkDecoder4x16Output(decoder, expected(12, 16), t)

	decoder.Update(true, true, false, true)
	checkDecoder4x16Output(decoder, expected(13, 16), t)

	decoder.Update(true, true, true, false)
	checkDecoder4x16Output(decoder, expected(14, 16), t)

	decoder.Update(true, true, true, true)
	checkDecoder4x16Output(decoder, expected(15, 16), t)
}

func checkDecoder2x4Output(decoder *Decoder2x4, expected []bool, t *testing.T) {
	for i, _ := range decoder.outputs {
		if decoder.outputs[i].Get() != expected[i] {
			t.FailNow()
		}
	}
}

func checkDecoder4x16Output(decoder *Decoder4x16, expected []bool, t *testing.T) {
	for i, _ := range decoder.outputs {
		if decoder.outputs[i].Get() != expected[i] {
			t.FailNow()
		}
	}
}

func checkDecoder3x8Output(decoder *Decoder3x8, expected []bool, t *testing.T) {
	for i, _ := range decoder.outputs {
		if decoder.outputs[i].Get() != expected[i] {
			t.FailNow()
		}
	}
}
