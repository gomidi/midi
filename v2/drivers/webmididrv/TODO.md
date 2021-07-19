
navigator.requestMIDIAccess()
    .then(onMIDISuccess, onMIDIFailure);

function onMIDISuccess(midiAccess) {
    console.log(midiAccess);

    var inputs = midiAccess.inputs;
    var outputs = midiAccess.outputs;
    
    for (var input of midiAccess.inputs.values())
        input.onmidimessage = getMIDIMessage;
    }
    
    inputs.forEach((midiInput) => {
	  // Do something with the MIDI input device
	});
	
	// Iterate through each connected MIDI output device
	outputs.forEach((midioutput) => {
	  // Do something with the MIDI output device 
	});
	
	midiInput.addEventListener('midimessage', (event) =&gt; {
	  // the `event` object will have a `data` property
	  // that contains an array of 3 numbers. For examples:
	  // [144, 63, 127]
	})
	
	outputsend([144, 63, 127]);
}

function onMIDIFailure() {
    console.log('Could not access your MIDI devices.');
}

function getMIDIMessage(message) {
    var command = message.data[0];
    var note = message.data[1];
    var velocity = (message.data.length > 2) ? message.data[2] : 0; // a velocity value might not be included with a noteOff command

    switch (command) {
        case 144: // noteOn
            if (velocity > 0) {
                noteOn(note, velocity);
            } else {
                noteOff(note);
            }
            break;
        case 128: // noteOff
            noteOff(note);
            break;
        // we could easily expand this switch statement to cover other types of commands such as controllers or sysex
    }
}

func jsonWrapper() js.Func {  
        jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
                if len(args) != 1 {
                        return "Invalid no of arguments passed"
                }
                inputJSON := args[0].String()
                fmt.Printf("input %s\n", inputJSON)
                pretty, err := prettyJson(inputJSON)
                if err != nil {
                        fmt.Printf("unable to convert to json %s\n", err)
                        return err.Error()
                }
                return pretty
        })
        return jsonFunc
}

func main() {  
        fmt.Println("Go Web Assembly")
        js.Global().Set("formatJSON", jsonWrapper())
        <-make(chan bool)
}

jsDoc := js.Global().Get("document")
        if !jsDoc.Truthy() {
            return "Unable to get document object"
        }
        jsonOuputTextArea := jsDoc.Call("getElementById", "jsonoutput")
        if !jsonOuputTextArea.Truthy() {
            return "Unable to get output text area"
        }
        inputJSON := args[0].String()




/*
For all other browsers that don’t support it natively, Chris Wilson’s WebMIDIAPIShim library (https://github.com/cwilso/WebMIDIAPIShim) is a polyfill for the Web MIDI API, of which Chris is a co-author. Simply including the shim script on your page will enable everything we’ve covered so far.

<script src="WebMIDIAPI.min.js"></script>
<script>
if (navigator.requestMIDIAccess) { //... returns true
</script>

This shim also requires Jazz-Soft.net’s Jazz-Plugin (https://jazz-soft.net/) to work, unfortunately, which means it’s an OK option for developers who want the flexibility to work in multiple browsers, but an extra barrier to mainstream adoption. Hopefully, within time, other browsers will adopt the Web MIDI API natively.
*/

/*
The MIDIMessageEvent object we get back contains a lot of information, but what we’re most interested in is the data array. This array typically contains three values (e.g. [144, 72, 64]). The first value tells us what type of command was sent, the second is the note value, and the third is velocity. The command type could be either “note on,” “note off,” controller (such as pitch bend or piano pedal), or some other kind of system exclusive (“sysex”) event unique to that device/manufacturer.

    A command value of 144 signifies a “note on” event, and 128 typically signifies a “note off” event.
    Note values are on a range from 0–127, lowest to highest. For example, the lowest note on an 88-key piano has a value of 21, and the highest note is 108. A “middle C” is 60.
    Velocity values are also given on a range from 0–127 (softest to loudest). The softest possible “note on” velocity is 1.
    A velocity of 0 is sometimes used in conjunction with a command value of 144 (which typically represents “note on”) to indicate a “note off” message, so it’s helpful to check if the given velocity is 0 as an alternate way of interpreting a “note off” message.
*/

