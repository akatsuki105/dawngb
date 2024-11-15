#ifndef LIBRETRO_H__
#define LIBRETRO_H__

#include <limits.h>
#include <stddef.h>
#include <stdint.h>

#ifndef __cplusplus
#if defined(_MSC_VER) && _MSC_VER < 1800 && !defined(SN_TARGET_PS3)
/* Hack applied for MSVC when compiling in C89 mode
 * as it isn't C99-compliant. */
#define bool unsigned char
#define true 1
#define false 0
#else
#include <stdbool.h>
#endif
#endif

#ifndef RETRO_CALLCONV
#if defined(__GNUC__) && defined(__i386__) && !defined(__x86_64__)
#define RETRO_CALLCONV __attribute__((cdecl))
#elif defined(_MSC_VER) && defined(_M_X86) && !defined(_M_X64)
#define RETRO_CALLCONV __cdecl
#else
#define RETRO_CALLCONV /* all other platforms only have one calling convention each */
#endif
#endif

#define RETRO_DEVICE_NONE 0
#define RETRO_DEVICE_JOYPAD 1

#define RETRO_DEVICE_ID_JOYPAD_B 0
#define RETRO_DEVICE_ID_JOYPAD_Y 1
#define RETRO_DEVICE_ID_JOYPAD_SELECT 2
#define RETRO_DEVICE_ID_JOYPAD_START 3
#define RETRO_DEVICE_ID_JOYPAD_UP 4
#define RETRO_DEVICE_ID_JOYPAD_DOWN 5
#define RETRO_DEVICE_ID_JOYPAD_LEFT 6
#define RETRO_DEVICE_ID_JOYPAD_RIGHT 7
#define RETRO_DEVICE_ID_JOYPAD_A 8
#define RETRO_DEVICE_ID_JOYPAD_X 9
#define RETRO_DEVICE_ID_JOYPAD_L 10
#define RETRO_DEVICE_ID_JOYPAD_R 11
#define RETRO_DEVICE_ID_JOYPAD_L2 12
#define RETRO_DEVICE_ID_JOYPAD_R2 13
#define RETRO_DEVICE_ID_JOYPAD_L3 14
#define RETRO_DEVICE_ID_JOYPAD_R3 15

#define RETRO_DEVICE_ID_JOYPAD_MASK 256

// retro_get_region() の戻り値
#define RETRO_REGION_NTSC 0
#define RETRO_REGION_PAL 1

#define RETRO_MEMORY_SAVE_RAM 0  // セーブデータ

#define RETRO_ENVIRONMENT_GET_SYSTEM_DIRECTORY 9
#define RETRO_ENVIRONMENT_SET_PIXEL_FORMAT 10
#define RETRO_ENVIRONMENT_GET_SAVE_DIRECTORY 31
#define RETRO_ENVIRONMENT_GET_INPUT_BITMASKS (51 | RETRO_ENVIRONMENT_EXPERIMENTAL)
#define RETRO_ENVIRONMENT_EXPERIMENTAL 0x10000

enum retro_pixel_format {
  /* 0RGB1555, native endian.
   * 0 bit must be set to 0.
   * This pixel format is default for compatibility concerns only.
   * If a 15/16-bit pixel format is desired, consider using RGB565. */
  RETRO_PIXEL_FORMAT_0RGB1555 = 0,

  /* XRGB8888, native endian.
   * X bits are ignored. */
  RETRO_PIXEL_FORMAT_XRGB8888 = 1,

  /* RGB565, native endian.
   * This pixel format is the recommended format to use if a 15/16-bit
   * format is desired as it is the pixel format that is typically
   * available on a wide range of low-power devices.
   *
   * It is also natively supported in APIs like OpenGL ES. */
  RETRO_PIXEL_FORMAT_RGB565 = 2,

  /* Ensure sizeof() == sizeof(int). */
  RETRO_PIXEL_FORMAT_UNKNOWN = INT_MAX
};

struct retro_system_info {
  /* All pointers are owned by libretro implementation, and pointers must
   * remain valid until it is unloaded. */

  const char *library_name;    /* Descriptive name of library. Should not
                                * contain any version numbers, etc. */
  const char *library_version; /* Descriptive version of core. */

  const char *valid_extensions; /* A string listing probably content
                                 * extensions the core will be able to
                                 * load, separated with pipe.
                                 * I.e. "bin|rom|iso".
                                 * Typically used for a GUI to filter
                                 * out extensions. */

  /* Libretro cores that need to have direct access to their content
   * files, including cores which use the path of the content files to
   * determine the paths of other files, should set need_fullpath to true.
   *
   * Cores should strive for setting need_fullpath to false,
   * as it allows the frontend to perform patching, etc.
   *
   * If need_fullpath is true and retro_load_game() is called:
   *    - retro_game_info::path is guaranteed to have a valid path
   *    - retro_game_info::data and retro_game_info::size are invalid
   *
   * If need_fullpath is false and retro_load_game() is called:
   *    - retro_game_info::path may be NULL
   *    - retro_game_info::data and retro_game_info::size are guaranteed
   *      to be valid
   *
   * See also:
   *    - RETRO_ENVIRONMENT_GET_SYSTEM_DIRECTORY
   *    - RETRO_ENVIRONMENT_GET_SAVE_DIRECTORY
   */
  bool need_fullpath;

  /* If true, the frontend is not allowed to extract any archives before
   * loading the real content.
   * Necessary for certain libretro implementations that load games
   * from zipped archives. */
  bool block_extract;
};

struct retro_game_geometry {
  unsigned base_width;  /* Nominal video width of game. */
  unsigned base_height; /* Nominal video height of game. */
  unsigned max_width;   /* Maximum possible width of game. */
  unsigned max_height;  /* Maximum possible height of game. */

  float aspect_ratio; /* Nominal aspect ratio of game. If
                       * aspect_ratio is <= 0.0, an aspect ratio
                       * of base_width / base_height is assumed.
                       * A frontend could override this setting,
                       * if desired. */
};

struct retro_system_timing {
  double fps;         /* FPS of video content. */
  double sample_rate; /* Sampling rate of audio. */
};

struct retro_system_av_info {
  struct retro_game_geometry geometry;
  struct retro_system_timing timing;
};

struct retro_game_info {
  const char *path; /* Path to game, UTF-8 encoded.
                     * Sometimes used as a reference for building other paths.
                     * May be NULL if game was loaded from stdin or similar,
                     * but in this case some cores will be unable to load `data`.
                     * So, it is preferable to fabricate something here instead
                     * of passing NULL, which will help more cores to succeed.
                     * retro_system_info::need_fullpath requires
                     * that this path is valid. */
  const void *data; /* Memory buffer of loaded game. Will be NULL
                     * if need_fullpath was set. */
  size_t size;      /* Size of memory buffer. */
  const char *meta; /* String of implementation specific meta-data. */
};

/* Callbacks */

/* Environment callback. Gives implementations a way of performing
 * uncommon tasks. Extensible. */
typedef bool(RETRO_CALLCONV *retro_environment_t)(unsigned cmd, void *data);
static retro_environment_t environ_cb;
static void _retro_set_environment(retro_environment_t cb) { environ_cb = cb; }
static bool call_environ_cb(unsigned cmd, void *data) { return environ_cb(cmd, data); }

/* Render a frame. Pixel format is 15-bit 0RGB1555 native endian
 * unless changed (see RETRO_ENVIRONMENT_SET_PIXEL_FORMAT).
 *
 * Width and height specify dimensions of buffer.
 * Pitch specifices length in bytes between two lines in buffer.
 *
 * For performance reasons, it is highly recommended to have a frame
 * that is packed in memory, i.e. pitch == width * byte_per_pixel.
 * Certain graphic APIs, such as OpenGL ES, do not like textures
 * that are not packed in memory.
 */
typedef void(RETRO_CALLCONV *retro_video_refresh_t)(const void *data, unsigned width, unsigned height, size_t pitch);
static retro_video_refresh_t video_cb;
static void _retro_set_video_refresh(retro_video_refresh_t cb) { video_cb = cb; }
static void call_video_cb(const void *data, unsigned width, unsigned height, size_t pitch) { video_cb(data, width, height, pitch); }

/* Renders a single audio frame. Should only be used if implementation
 * generates a single sample at a time.
 * Format is signed 16-bit native endian.
 */
typedef void(RETRO_CALLCONV *retro_audio_sample_t)(int16_t left, int16_t right);
static retro_audio_sample_t audio_cb;
static void _retro_set_audio_sample(retro_audio_sample_t cb) { audio_cb = cb; }
static void call_audio_cb(int16_t left, int16_t right) { audio_cb(left, right); }

/* Renders multiple audio frames in one go.
 *
 * One frame is defined as a sample of left and right channels, interleaved.
 * I.e. int16_t buf[4] = { l, r, l, r }; would be 2 frames.
 * Only one of the audio callbacks must ever be used.
 */
typedef size_t(RETRO_CALLCONV *retro_audio_sample_batch_t)(const int16_t *data, size_t frames);
static retro_audio_sample_batch_t audio_batch_cb;
static void _retro_set_audio_sample_batch(retro_audio_sample_batch_t cb) { audio_batch_cb = cb; }
static void call_audio_batch_cb(const int16_t *data, size_t frames) { audio_batch_cb(data, frames); }

/* Polls input. */
typedef void(RETRO_CALLCONV *retro_input_poll_t)(void);
static retro_input_poll_t input_poll_cb;
static void _retro_set_input_poll(retro_input_poll_t cb) { input_poll_cb = cb; }
static void call_input_poll_cb(void) { input_poll_cb(); }

/* Queries for input for player 'port'. device will be masked with
 * RETRO_DEVICE_MASK.
 *
 * Specialization of devices such as RETRO_DEVICE_JOYPAD_MULTITAP that
 * have been set with retro_set_controller_port_device()
 * will still use the higher level RETRO_DEVICE_JOYPAD to request input.
 */
typedef int16_t(RETRO_CALLCONV *retro_input_state_t)(unsigned port, unsigned device, unsigned index, unsigned id);
static retro_input_state_t input_state_cb;
static void _retro_set_input_state(retro_input_state_t cb) { input_state_cb = cb; }
static int16_t call_input_state_cb(unsigned port, unsigned device, unsigned index, unsigned id) { return input_state_cb(port, device, index, id); }

#endif
