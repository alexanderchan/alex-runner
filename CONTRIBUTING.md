## How to record and playback tapes

https://github.com/charmbracelet/vhs?tab=readme-ov-file#record-tapes


1. record

```bash
vhs record > cassette.tape
```

2. Set the output file settings in the tape

```bash
Output ./public/images/demo.gif
Set TypingSpeed 0.1s
```

3. Generate the file

```bash
	vhs cassette.tape
```

There's some make scripts to help with this.  Should make an examples directory with package.json for the recording next time.
