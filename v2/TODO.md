# sysex: implement

SysEx commands carried on a MIDI 1.0 DIN cable may use the "dropped
   0xF7" construction [MIDI].  In this coding method, the 0xF7 octet is
   dropped from the end of the SysEx command, and the status octet of
   the next MIDI command acts both to terminate the SysEx command and
   start the next command. 

Inbetween these two status bytes, any number of data bytes (all having bit #7 clear, ie, 0 to 127 value) may be sent.

The purpose of the remaining data bytes, however many there may be, are determined by the manufacturer of a product. Typically, manufacturers follow the Manufacturer ID with a Model Number ID byte so that a device can not only determine that it's got a SysEx message for the correct manufacturer, but also has a SysEx message specifically for this model. Then, after the Model ID may be another byte that the device uses to determine what the SysEx message is supposed to be for, and therefore, how many more data bytes follow. Some manufacturers have a checksum byte, (usually right before the 0xF7) which is used to check the integrity of the message's transmission. 

The 0xF7 Status byte is dedicated to marking the end of a SysEx message. It should never occur without a preceding 0xF0 Status. In the event that a device experiences such a condition (ie, maybe the MIDI cable was connected during the transmission of a SysEx message), the device should ignore the 0xF7.

Furthermore, although the 0xF7 is supposed to mark the end of a SysEx message, in fact, any status (except for Realtime Category messages) will cause a SysEx message to be considered "done" (ie, actually "aborted" is a better description since such a scenario indicates an abnormal MIDI condition). For example, if a 0x90 happened to be sent sometime after a 0xF0 (but before the 0xF7), then the SysEx message would be considered aborted at that point. It should be noted that, like all System Common messages, SysEx cancels any current running status. In other words, the next Voice Category message (after the SysEx message) must begin with a Status. 

# undefined syscommon and undefined Real-Time commands: implement

 [MIDI] reserves the undefined System Common commands 0xF4 and 0xF5
   and the undefined System Real-Time commands 0xF9 and 0xFD for future
   use. 


# make tests for midi channel messages



# make transparent running status to explicit status reader; make it the default in listener, let it start listening at the first explicit status

# make transparent running status writer

# here is the question if it is not better to have some midi.Buffer that tracks status bytes for reading and writing.
such a buffer could be used to convert to explicit status codes (which would be done inside a smf reader and a driver in port)
or to compress with the help of running status (which would be done inside a smf writer and a driver out port)

we need an interface with optional methods (=callbacks) for sysex messages (with are buffered by the drivers),
realtime messages (with are send immediatly) and system common messages (which are also send immediatly).
the default method receives channel messages (or could possibly receive any kind of message)

# improve sysex

# Pipelines / Builders

# Test SMPTE in smf

# Test midi clock etc. realtime and syscommon messages

# Test sysex

