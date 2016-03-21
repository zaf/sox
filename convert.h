/*
	Sound file conversion using soxlib

	Copyright (C) 2016, Lefteris Zafiris <zaf.000@gmail.com>

	This program is free software, distributed under the terms of
	the BSD 3-Clause License. See the LICENSE file
	at the top of the source tree.

*/

#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <malloc.h>
#include <sox.h>

typedef struct snd_file {
	void *buff;
	size_t size;
} snd_file;

snd_file convert_snd(snd_file *in_snd, char* in_format, char* out_format) {
	snd_file out_snd = { NULL, 0 };
	static sox_format_t *in, *out;
	sox_effects_chain_t *chain;
	sox_effect_t *e;
	char *args[10];
	char *ftype;
	sox_signalinfo_t interm_signal, out_signal;
	sox_encodinginfo_t out_encoding;

	if (strcmp(out_format, "ulaw") == 0 || strcmp(out_format, "ul") == 0) {
		ftype = "ul";
		out_signal = (sox_signalinfo_t) { 8000, 1, 0, 0, NULL };
		out_encoding = (sox_encodinginfo_t) {
			SOX_ENCODING_ULAW,
			8,
			0,
			sox_option_default,
			sox_option_default,
			sox_option_default,
			sox_false
		};
	} else if (strcmp(out_format, "alaw") == 0 || strcmp(out_format, "al") == 0) {
		ftype = "al";
		out_signal = (sox_signalinfo_t) { 8000, 1, 0, 0, NULL };
		out_encoding = (sox_encodinginfo_t) {
			SOX_ENCODING_ALAW,
			8,
			0,
			sox_option_default,
			sox_option_default,
			sox_option_default,
			sox_false
		};
	} else if (strcmp(out_format, "gsm") == 0) {
		ftype = "gsm";
		out_signal = (sox_signalinfo_t) { 8000, 1, 0, 0, NULL };
		out_encoding = (sox_encodinginfo_t) {
			SOX_ENCODING_GSM,
			0,
			1,
			sox_option_default,
			sox_option_default,
			sox_option_default,
			sox_false
		};
	} else if (strcmp(out_format, "wav") == 0 || strcmp(out_format, "wav16") == 0) {
		ftype = "wav";
		out_signal = (sox_signalinfo_t) { 16000, 1, 0, 0, NULL };
		out_encoding = (sox_encodinginfo_t) {
			SOX_ENCODING_SIGN2,
			16,
			0,
			sox_option_default,
			sox_option_default,
			sox_option_default,
			sox_false
		};
	} else if (strcmp(out_format, "wav8") == 0) {
		ftype = "wav";
		out_signal = (sox_signalinfo_t) { 8000, 1, 0, 0, NULL };
		out_encoding = (sox_encodinginfo_t) {
			SOX_ENCODING_SIGN2,
			16,
			0,
			sox_option_default,
			sox_option_default,
			sox_option_default,
			sox_false
		};
	} else if (strcmp(out_format, "mp3") == 0) {
		ftype = "mp3";
		out_signal = (sox_signalinfo_t) { 16000, 1, 0, 0, NULL };
		out_encoding = (sox_encodinginfo_t) {
			SOX_ENCODING_MP3,
			0,
			32,
			sox_option_default,
			sox_option_default,
			sox_option_default,
			sox_false
		};
	} else if (strcmp(out_format, "ogg") == 0) {
		ftype = "ogg";
		out_signal = (sox_signalinfo_t) { 16000, 1, 0, 0, NULL };
		out_encoding = (sox_encodinginfo_t) {
			SOX_ENCODING_VORBIS,
			0,
			4,
			sox_option_default,
			sox_option_default,
			sox_option_default,
			sox_false
		};
	} else if (strcmp(out_format, "flac") == 0) {
		ftype = "flac";
		out_signal = (sox_signalinfo_t) { 16000, 1, 0, 0, NULL };
		out_encoding = (sox_encodinginfo_t) {
			SOX_ENCODING_FLAC,
			0,
			6,
			sox_option_default,
			sox_option_default,
			sox_option_default,
			sox_false
		};
	} else {
		return out_snd;
	}

	in = sox_open_mem_read(in_snd->buff, in_snd->size, NULL, NULL, in_format);
	if (in == NULL || in->encoding.encoding == SOX_ENCODING_UNKNOWN ) {
		return out_snd;
	}
	out = sox_open_memstream_write((char**)&out_snd.buff, &out_snd.size, &out_signal, &out_encoding, ftype, NULL);
	if (out == NULL) {
		sox_close(in);
		return out_snd;
	}

	chain = sox_create_effects_chain(&in->encoding, &out->encoding);
	interm_signal = in->signal;
	e = sox_create_effect(sox_find_effect("input"));
	args[0] = (char *) in;
	if ( sox_effect_options(e, 1, args) != SOX_SUCCESS ||
			sox_add_effect(chain, e, &interm_signal, &in->signal) != SOX_SUCCESS ) {
		goto END;
	}
	free(e);

	if (in->signal.rate != out->signal.rate) {
		e = sox_create_effect(sox_find_effect("rate"));
		if ( sox_effect_options(e, 0, NULL) != SOX_SUCCESS ||
				sox_add_effect(chain, e, &interm_signal, &out->signal) != SOX_SUCCESS ) {
			goto END;
		}
		free(e);
	}

	if (in->signal.channels != out->signal.channels) {
		e = sox_create_effect(sox_find_effect("channels"));
		if ( sox_effect_options(e, 0, NULL) != SOX_SUCCESS ||
				sox_add_effect(chain, e, &interm_signal, &out->signal) != SOX_SUCCESS ) {
			goto END;
		}
		free(e);
	}

	e = sox_create_effect(sox_find_effect("output"));
	args[0] = (char *) out;
	if ( sox_effect_options(e, 1, args) != SOX_SUCCESS ||
			sox_add_effect(chain, e, &interm_signal, &out->signal) != SOX_SUCCESS ) {
		goto END;
	}

	sox_flow_effects(chain, NULL, NULL);
END:
	free(e);
	sox_delete_effects_chain(chain);
	sox_close(in);
	sox_close(out);

	return out_snd;
}
