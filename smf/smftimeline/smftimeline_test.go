package smftimeline

import (
	"testing"

	"gitlab.com/gomidi/midi/smf"
)

func TestForwardNBars(t *testing.T) {

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

		tl.ForwardNBars(test.nbars)

		if tl.cursor != int64(test.expected) {
			t.Errorf("[%v] expecting %v, got: %v", i, test.expected, tl.cursor)
		}

	}

}

func TestForward(t *testing.T) {
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
			tl.Forward(step[0], step[1])
		}

		if tl.cursor != test.expected {
			t.Errorf("[%v] expecting %v, got: %v", i, test.expected, tl.cursor)
		}

	}

}
