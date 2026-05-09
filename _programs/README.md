These are a few example programs I have written.

All of them  were generated from [cmd/generator](../cmd/generator) so utilise some shared routines/code.

You can regenerate them by running `make`

# ascii.bin

![ascii.bin](screenshots/ascii.png?raw=true "ascii.bin")

Renders the ASCII table from 32 - 126. 

Note the letters don't have lowercase forms as I couldn't be bothered to craft them.

# brush.bin

![brush.bin](screenshots/brush.png?raw=true "brush.bin")

Use the arrow keys to move a 'snake' style brush around the screen.

# me.bin

![me.bin](screenshots/me.png?raw=true "me.bin")

Render some stuff about me

# text-writer.bin

![text-writer.bin](screenshots/text-writer.png?raw=true "text-writer.bin")

Use the keyboard to type ASCII characters, which get rendered on the display.
Hit enter to perform a carriage return.

Note: not all keys work or will render garbage, e.g. backspace.
Also modifier keys are not supported so rendering symbols that require them (e.g. shift) won't work.

# counter.bin

Displays a 16-bit hexadecimal counter that increments continuously from `0x0000` to `0xFFFF`, then wraps back to zero.

# sine.bin

Displays a scrolling sine wave animation. Thirty dots are plotted across the display — one per character column — using a precomputed lookup table. Each frame the table rotates by one entry, making the wave appear to scroll horizontally.