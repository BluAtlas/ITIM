# iPhone Transparent Icon Maker

This tool will take a background and a set of partially transparent icon images, scale the background to fit a desired iPhone, and then create new icon images matching the background with the transparency to be used with the shortcuts/bookmarks function to create the illusion of transparent icons on the home screen.

## How to use

On a first run, ITIM will create two folders, input and output. Place a "background.jpg/png" in the input folder, and any partially transparent icon files in "input/icons". Icon files should be a 1:1 aspect ratio, and the background will be scaled inwards as required to fit the screen. These icon files are sorted by file name, so rename them however you please to sort their locations on the home screen. The last few icons in the folder will be put on the dock.

Then run one of the following commands:

    for iPhone 12 Pro Max layout:
    ./ITIM 1 ['png' or 'jpg']

    for iPhone 12 mini layout:
    ./ITIM 2 ['png' or 'jpg']

    To create icons for a custom sizing:
    ./ITIM [FILETYPE PHONESIZEX PHONESIZEY APPSIZE FROMTOP FROMLEFT BETWEENX BETWEENY FROMBOTTOM FROMLEFTDOCK MAXAPPS DOCKCOUNT]

    FILETYPE is either 'png'or 'jpg'
    PHONESIZEX and PHONESIZEY should match the resolution of a screenshot from the desired phone
    APPSIZE should be width in pixels of an app icon
    FROMTOP should be the amount of pixels between the top of the image, and the top of the top-left most app
    FROMLEFT should be the amount of pixels between the left side of the image and the left side of the top-left most app
    BETWEENX should be the amount of pixels between two app icons in the left-right direction
    BETWEENY should be the amount of pixels between two app icons in the up-down direction
    FROMBOTTOM should be the amount of pixels between the bottom of the image and the left most app on the dock
    FROMLEFTDOCK should be the amount of pixels between the left side of the image and the left most app on the dock
    MAXAPPS should be the maximum number of apps that can be shown on screen, not including the dock
    DOCKCOUNT should be the amount of apps on your dock