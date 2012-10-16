This tool is designed to remove the most detailed zoom levels from a [Minecraft Overviewer] map.
It takes care of both the removal of the rendered files as well as the patching the `overviewerConfig.js`.

# Usage
You have to be in the folder containing the `overviewerConfig.js`

    $ zoomreduce remove -n 2 -w world -w nether -w the_end

This will remove the 2 highest zoom leves of `world`, `nether` and `the_end`.

[Minecraft Overviewer]: https://github.com/overviewer/Minecraft-Overviewer
---
Version 0.1
