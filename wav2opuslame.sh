#! /bin/bash

# Check config file integrity
if [ ! -e ./etc/config.conf ]; then
    >&2 echo "CRITICAL FAILURE: No config file found. Get a fresh one from ./etc/template.conf (unless you deleted that too)"
    exit 1
fi
source ./etc/config.conf

# Variable check
if [ -z "$ARTIST" ]; then
    >&2 echo "ERROR: ARTIST tag missing in config file"
    ERR_STATUS=1
fi

if [ -z "$TITLE" ]; then
    >&2 echo "ERROR: TITLE tag missing in config file"
    ERR_STATUS=1
fi

if [ -z "$ALBUM" ]; then 
    >&2 echo "ERROR: ALBUM tag missing in config file"
    ERR_STATUS=1
fi

if [ -z "$LAME_BITRATE" ]; then
    >&2 echo "ERROR: LAME BITRATE variable missing in config file"
    ERR_STATUS=1
fi

if [ -z "$OPUS_BITRATE" ]; then
    >&2 echo "ERROR: OPUS BITRATE tag missing in config file"
    ERR_STATUS=1
fi

# Check if input file is readable and output directory is writeable
if test ! -r "$REC_DIR/$1" || test ! -w "$OUT_DIR/$2"; then
    echo "I do not have the proper file permissions for input/output files"
    ERR_STATUS=1
fi

if [ "$ERR_STATUS" ]; then
    >&2 echo "Errors found. See output above for details."
    exit 1
fi

echo "Config file integrity check passed.\n"

# DEBUG OUTPUT: Will make this a toggable command later
echo "ARTIST:$ARTIST"
echo "TITLE:$TITLE"
echo "ALBUM:$ALBUM"
echo "REC_DIR:$REC_DIR"
echo "OUT_DIR:$OUT_DIR"
echo "LAME_BITRATE:$LAME_BITRATE"
echo "OPUS_BITRATE:$OPUS_BITRATE"



# TODO: CHECK OPUSENC AND LAME EXECUTABLES





# TODO: FIRE ENCODING COMMANDS