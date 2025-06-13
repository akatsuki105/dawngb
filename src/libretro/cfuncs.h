#ifndef DAWNGB_CFUNCS
#define DAWNGB_CFUNCS

static retro_environment_t environ_cb;
static void _retro_set_environment(retro_environment_t cb) { environ_cb = cb; }
static bool call_environ_cb(unsigned cmd, void *data) { return environ_cb(cmd, data); }

static struct retro_log_callback logging;
static retro_log_printf_t log_cb;

static retro_video_refresh_t video_cb;
static void _retro_set_video_refresh(retro_video_refresh_t cb) { video_cb = cb; }
static void call_video_cb(const void *data, unsigned width, unsigned height, size_t pitch) { video_cb(data, width, height, pitch); }

static retro_audio_sample_t audio_cb;
static void _retro_set_audio_sample(retro_audio_sample_t cb) { audio_cb = cb; }
static void call_audio_cb(int16_t left, int16_t right) { audio_cb(left, right); }

static retro_audio_sample_batch_t audio_batch_cb;
static void _retro_set_audio_sample_batch(retro_audio_sample_batch_t cb) { audio_batch_cb = cb; }
static void call_audio_batch_cb(const int16_t *data, size_t frames) { audio_batch_cb(data, frames); }

static retro_input_poll_t input_poll_cb;
static void _retro_set_input_poll(retro_input_poll_t cb) { input_poll_cb = cb; }
static void call_input_poll_cb(void) { input_poll_cb(); }

static retro_input_state_t input_state_cb;
static void _retro_set_input_state(retro_input_state_t cb) { input_state_cb = cb; }
static int16_t call_input_state_cb(unsigned port, unsigned device, unsigned index, unsigned id) { return input_state_cb(port, device, index, id); }

#endif  // DAWNGB_CFUNCS
