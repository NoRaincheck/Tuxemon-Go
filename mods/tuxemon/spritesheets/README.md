To create the spritesheet, the following code was used from the root `tuxemon` repo

```py
# /// script
# requires-python = ">=3.10"
# dependencies = [
#     "pillow",
# ]
# ///

from PIL import Image

sprites = [
    [
        "mods/tuxemon/sprites/adventurer_front.png",
        "mods/tuxemon/sprites/adventurer_front_walk.000.png",
        "mods/tuxemon/sprites/adventurer_front_walk.001.png",
    ],
    [
        "mods/tuxemon/sprites/adventurer_back.png",
        "mods/tuxemon/sprites/adventurer_back_walk.000.png",
        "mods/tuxemon/sprites/adventurer_back_walk.001.png",
    ],
    [
        "mods/tuxemon/sprites/adventurer_left.png",
        "mods/tuxemon/sprites/adventurer_left_walk.000.png",
        "mods/tuxemon/sprites/adventurer_left_walk.001.png",
    ],
    [
        "mods/tuxemon/sprites/adventurer_right.png",
        "mods/tuxemon/sprites/adventurer_right_walk.000.png",
        "mods/tuxemon/sprites/adventurer_right_walk.001.png",
    ],
]

# get the max width and height of the sprites
images = []
max_total_width = 0
max_total_height = 0
for sprite_row in sprites:
    images = [Image.open(sprite) for sprite in sprite_row]
    total_width = sum(image.width for image in images)
    max_height = max(image.height for image in images)

    max_total_width = max(max_total_width, total_width)
    max_total_height += max_height

# create a new image with the max width and height
new_image = Image.new("RGBA", (max_total_width, max_total_height))

# paste the sprites into the new image
for height_indx, sprite_row in enumerate(sprites):
    for width_indx, sprite in enumerate(sprite_row):
        image = Image.open(sprite)
        new_image.paste(image, (width_indx * image.width, height_indx * image.height))

new_image.save("mods/tuxemon/spritesheets/adventurer_walk.png")
```
