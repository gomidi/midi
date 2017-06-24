package meta

/*
http://www.somascape.org/midi/tech/mfile.html

SMPTE Offset

FF 54 05 hr mn se fr ff
hr is a byte specifying the hour, which is also encoded with the SMPTE format (frame rate), just as it is in MIDI Time Code, i.e. 0rrhhhhh, where :
rr = frame rate : 00 = 24 fps, 01 = 25 fps, 10 = 30 fps (drop frame), 11 = 30 fps (non-drop frame)
hhhhh = hour (0-23)
mn se are 2 bytes specifying the minutes (0-59) and seconds (0-59), respectively.
fr is a byte specifying the number of frames (0-23/24/28/29, depending on the frame rate specified in the hr byte).
ff is a byte specifying the number of fractional frames, in 100ths of a frame (even in SMPTE-based tracks using a different frame subdivision, defined in the MThd chunk).

This optional event, if present, should occur at the start of a track, at time = 0, and prior to any MIDI events. It is used to specify the SMPTE time at which the track is to start.

For a format 1 MIDI file, a SMPTE Offset Meta event should only occur within the first MTrk chunk.

*/

/* http://www.somascape.org/midi/tech/mfile.html
Sequencer Specific Event

FF 7F length data

The first 1 or 3 bytes of data is a manufacturer's ID code (same format as for System Exclusive messages). This optional event can be used to store sequencer-specific information.

*/

/* http://www.somascape.org/midi/tech/mfile.html
Program Name

FF 08 length text

This optional event is used to embed the patch/program name that is called up by the immediately subsequent Bank Select and Program Change messages. It serves to aid the end user in making an intelligent program choice when using different hardware.

This event may appear anywhere in a track, and there may be multiple occurrences within a track.
*/
