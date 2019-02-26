package smftimeline

import (
	"reflect"
	"testing"

	"gitlab.com/gomidi/midi/smf"
)

func TestPlanned(t *testing.T) {

	ticks := smf.MetricTicks(960)

	var tests = []struct {
		start     int64
		until     int64
		nbars     uint32
		steps     [2]uint32
		deltas    []int32
		times     int
		remaining int
	}{
		{

			0, int64(ticks.Ticks4th() * 8),
			1, [2]uint32{0, 0},
			[]int32{int32(ticks.Ticks4th() * 4)},
			1,
			0,
		},
		{

			0, int64(ticks.Ticks4th() * 8),
			1, [2]uint32{1, 4},
			[]int32{int32(ticks.Ticks4th() * 5)},
			1,
			0,
		},
		{

			0, int64(ticks.Ticks4th() * 8),
			1, [2]uint32{1, 4},
			[]int32{int32(ticks.Ticks4th() * 5), 0, 0},
			3,
			0,
		},
		{

			int64(ticks.Ticks4th() * 4), -1,
			2, [2]uint32{1, 4},
			[]int32{int32(ticks.Ticks4th() * 5), 0, 0},
			3,
			0,
		},
		{

			0, int64(ticks.Ticks4th() * 8),
			2, [2]uint32{1, 4},
			[]int32{},
			1,
			1,
		},
		{

			int64(ticks.Ticks4th() * 5), int64(ticks.Ticks4th() * 8),
			1, [2]uint32{0, 0},
			[]int32{},
			1,
			0,
		},
		{

			int64(ticks.Ticks4th() * 4), int64(ticks.Ticks4th() * 8),
			1, [2]uint32{0, 0},
			[]int32{0},
			1,
			0,
		},
		{

			int64(ticks.Ticks4th() * 1), int64(ticks.Ticks4th() * 8),
			1, [2]uint32{0, 0},
			[]int32{int32(ticks.Ticks4th() * 3)},
			1,
			0,
		},
	}

	for i, test := range tests {
		//		if i != 2 {
		//			continue
		//		}
		//		fmt.Printf("test %v\n", i)
		var tl TimeLine
		tl.ticks = ticks
		tl.Reset()

		var deltas []int32

		for j := 0; j < test.times; j++ {
			tl.Plan(test.nbars, test.steps[0], test.steps[1], func(_delta int32) {
				deltas = append(deltas, _delta)
			})
		}

		tl.cursor = test.start
		tl.runCallbacks(test.until)

		if !reflect.DeepEqual(deltas, test.deltas) && (len(deltas) != 0 || len(test.deltas) != 0) {
			t.Errorf("[%v] expecting deltas %v, got: %v", i, test.deltas, deltas)
		}

		if got, expected := len(tl.plannedCallbacks), test.remaining; got != expected {
			t.Errorf("[%v] expecting remaining %v, got: %v", i, test.remaining, got)
		}

	}

}

func TestForwardNBars(t *testing.T) {
	//	t.Skip()
	var tests = []struct {
		ticks    smf.MetricTicks
		start    uint32
		timeSigs [][3]int64
		nbars    uint32
		expected uint32
	}{
		{
			smf.MetricTicks(960),
			0,
			[][3]int64{}, // timeSigs
			0,            // nbars
			0,            // result
		}, {
			smf.MetricTicks(960),
			0,
			[][3]int64{},                        // timeSigs
			1,                                   // nbars
			smf.MetricTicks(960).Ticks4th() * 4, // result
		},
		{
			smf.MetricTicks(960),
			0,
			[][3]int64{},                            // timeSigs
			4,                                       // nbars
			smf.MetricTicks(960).Ticks4th() * 4 * 4, // result
		},
		{
			smf.MetricTicks(960),
			0,
			[][3]int64{[3]int64{int64(smf.MetricTicks(960).Ticks4th() * 4 * 2), 3, 4}}, // timeSigs
			5, // nbars
			smf.MetricTicks(960).Ticks4th()*4*2 + smf.MetricTicks(960).Ticks4th()*3*3, // result
		},
		{
			smf.MetricTicks(960),
			smf.MetricTicks(960).Ticks4th(),
			[][3]int64{[3]int64{int64(smf.MetricTicks(960).Ticks4th() * 4 * 2), 3, 4}}, // timeSigs
			5, // nbars
			smf.MetricTicks(960).Ticks4th()*4*2 + smf.MetricTicks(960).Ticks4th()*3*3, // result
		},
		{
			smf.MetricTicks(960),
			smf.MetricTicks(960).Ticks4th() * 9,
			[][3]int64{[3]int64{int64(smf.MetricTicks(960).Ticks4th() * 4 * 2), 3, 4}}, // timeSigs
			5, // nbars
			smf.MetricTicks(960).Ticks4th()*4*2 + smf.MetricTicks(960).Ticks4th()*3*5, // result
		},
		{
			smf.MetricTicks(960),
			smf.MetricTicks(960).Ticks4th() * 12,
			[][3]int64{[3]int64{int64(smf.MetricTicks(960).Ticks4th() * 4 * 2), 3, 4}}, // timeSigs
			5, // nbars
			smf.MetricTicks(960).Ticks4th()*4*2 + smf.MetricTicks(960).Ticks4th()*3*6, // result
		},
		{
			smf.MetricTicks(1920),
			smf.MetricTicks(1920).Ticks4th(),
			[][3]int64{[3]int64{int64(smf.MetricTicks(1920).Ticks4th() * 4 * 2), 3, 4}}, // timeSigs
			5, // nbars
			smf.MetricTicks(1920).Ticks4th()*4*2 + smf.MetricTicks(1920).Ticks4th()*3*3, // result
		},
		{
			smf.MetricTicks(1920),
			smf.MetricTicks(1920).Ticks4th() * 9,
			[][3]int64{[3]int64{int64(smf.MetricTicks(1920).Ticks4th() * 4 * 2), 3, 4}}, // timeSigs
			5, // nbars
			smf.MetricTicks(1920).Ticks4th()*4*2 + smf.MetricTicks(1920).Ticks4th()*3*5, // result
		},
		{
			smf.MetricTicks(1920),
			smf.MetricTicks(1920).Ticks4th() * 12,
			[][3]int64{[3]int64{int64(smf.MetricTicks(1920).Ticks4th() * 4 * 2), 3, 4}}, // timeSigs
			5, // nbars
			smf.MetricTicks(1920).Ticks4th()*4*2 + smf.MetricTicks(1920).Ticks4th()*3*6, // result
		},
	}

	for i, test := range tests {
		var tl TimeLine
		tl.ticks = test.ticks
		tl.timeSigs = test.timeSigs

		tl.Reset()
		tl.cursor = int64(test.start)

		tl.Forward(test.nbars, 0, 0)

		if tl.cursor != int64(test.expected) {
			t.Errorf("[%v] expecting %v, got: %v", i, test.expected, tl.cursor)
		}

	}

}

func TestForward(t *testing.T) {
	//	t.Skip()
	var tests = []struct {
		ticks    smf.MetricTicks
		steps    [][2]uint32
		expected int64
	}{
		{
			smf.MetricTicks(960),
			[][2]uint32{}, // steps
			int64(0),      // result
		}, {
			smf.MetricTicks(960),
			[][2]uint32{[2]uint32{4, 4}},               // steps
			int64(smf.MetricTicks(960).Ticks4th() * 4), // result
		},
		{
			smf.MetricTicks(960),
			[][2]uint32{{4, 4}, {6, 8}},                // steps
			int64(smf.MetricTicks(960).Ticks4th() * 7), // result
		},
	}

	for i, test := range tests {
		var tl TimeLine
		tl.ticks = test.ticks

		tl.Reset()

		for _, step := range test.steps {
			tl.Forward(0, step[0], step[1])
		}

		if tl.cursor != test.expected {
			t.Errorf("[%v] expecting %v, got: %v", i, test.expected, tl.cursor)
		}

	}

}
