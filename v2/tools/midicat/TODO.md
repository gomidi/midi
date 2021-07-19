
Add this list of commands:

- read: reads a track from an smf file and prints it to stdout
- write: writes midi track data from stdin to a smf (may append or create a new one)
- play: plays a smf file to stdout
- record: records stdin to a smf file
- send: sends live midi from stdin to a tcp receiver
- receive: receives midi via tcp and writes it to stdout
- filter: passes midi data from stdin to stdout while filtering it.
          filtering is possible
          - by channel (channel parameter)
          - by message kind (kind parameter), e.g. channelmsg, metamsg,sysexmsg etc.
          - by message type (type parameter), e.g. controlchange,noteon,noteoff
          - by value: max/min, e.g.  
             +controlchange(3)[12:120]
            or exact
             +controlchange(3)[0,127],-controlchange(7)[47:100]
        
          all params work like this:
            + prefix denotes passthrough
            - prefix denotes blocking
        