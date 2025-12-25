#ifndef DAWNGB_INPUT
#define DAWNGB_INPUT

#include "libretro.h"

static const struct retro_controller_description controllers[] = {
    {"Nintendo Game Boy", RETRO_DEVICE_SUBCLASS(RETRO_DEVICE_JOYPAD, 0)},
};

static const struct retro_controller_info ports[] = {
    {controllers, 1},
    {NULL, 0},
};

static struct retro_input_descriptor descriptors_1p[] = {
  { 0, RETRO_DEVICE_JOYPAD, 0, RETRO_DEVICE_ID_JOYPAD_LEFT,  "Left" },
  { 0, RETRO_DEVICE_JOYPAD, 0, RETRO_DEVICE_ID_JOYPAD_UP,    "Up" },
  { 0, RETRO_DEVICE_JOYPAD, 0, RETRO_DEVICE_ID_JOYPAD_DOWN,  "Down" },
  { 0, RETRO_DEVICE_JOYPAD, 0, RETRO_DEVICE_ID_JOYPAD_RIGHT, "Right" },
  { 0, RETRO_DEVICE_JOYPAD, 0, RETRO_DEVICE_ID_JOYPAD_B, "B" },
  { 0, RETRO_DEVICE_JOYPAD, 0, RETRO_DEVICE_ID_JOYPAD_A, "A" },
  { 0, RETRO_DEVICE_JOYPAD, 0, RETRO_DEVICE_ID_JOYPAD_SELECT, "Select" },
  { 0, RETRO_DEVICE_JOYPAD, 0, RETRO_DEVICE_ID_JOYPAD_START, "Start" },
  { 0 },
};

#endif  // DAWNGB_INPUT
