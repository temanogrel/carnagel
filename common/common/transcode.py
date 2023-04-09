import time
import json
import subprocess
import math
import os


class TranscodingError(RuntimeError):
    pass


class TranscodingUtility:
    def __init__(self, ffmpeg: str, ffprobe: str, vcs: str, watermark: str, quiet=True):

        self.vcs = vcs
        self.ffmpeg = ffmpeg
        self.ffprobe = ffprobe
        self.watermark = watermark
        self.quiet = quiet

    def _run_args(self, args: list):
        """
        Execute the given argument list

        :param args:
        :return:
        """

        try:

            if not self.quiet:
                return subprocess.check_call(args, stderr=subprocess.STDOUT)

            with open(os.devnull, 'w') as dev_null:
                return subprocess.check_call(args, stdout=dev_null, stderr=subprocess.STDOUT)

        except subprocess.CalledProcessError as e:
            raise TranscodingError(e)

    def generate_thumbnails(self, source: str, target: str):
        """
        Generates a bunch of thumbnails based on the duration of the video to the target path

        :param source:
        :param target:
        :return: Yields each new file
        """

        video, audio, data = self._probe_source(source)

        duration = float(data['format']['duration'])

        if duration < 10 * 60:
            generate_images = 16

        elif duration < 30 * 60:
            generate_images = 32

        elif duration < 2 * 60 * 60:
            generate_images = 64

        else:
            generate_images = 88

        path, ext = os.path.splitext(target)

        # create the folde rif it does not exist
        if not os.path.exists(os.path.dirname(target)):
            os.makedirs(os.path.dirname(target))

        # convert it into something that ffmpeg understands how to format
        ffmpeg_target = path + '_%02d' + ext

        args = [
            self.ffmpeg,

            # Overwrite if file exists
            '-y',

            '-loglevel', '24' if self.quiet else '32',

            '-i', source,

            '-vf', 'thumbnail=100',
            '-vf', 'fps=1/{}'.format(math.ceil(duration / (generate_images - 1))),

            '-q:v', '1',

            ffmpeg_target
        ]

        self._run_args(args)

        for x in range(1, generate_images + 1):
            image_path = path + '_{0:02d}'.format(x) + ext
            captured_at = math.floor((duration / generate_images) * x)

            yield captured_at, image_path

    def to_h264(self, source: str, target: str, crf: int, preset: str, tune: str, threads: int, should_watermark=True):
        """
        Convert a video to h264

        :param source:
        :param target:

        :return:
        """

        if not os.path.exists(source):
            raise FileNotFoundError(source)

        args = [
            self.ffmpeg,

            # Overwrite if file exists
            '-y',

            # Output level (Warnings)
            '-loglevel', '24' if self.quiet else '32',

            # Input filter
            '-i', source,

            # Video
            '-c:v', 'libx264', '-crf', str(crf),

            # Audio
            '-c:a', 'libfdk_aac', '-preset', str(preset),

            # Frame rate
            '-tune', str(tune),

            # Number of threads to use
            '-threads', str(threads),
        ]

        if should_watermark:
            args.append('-vf')
            args.append(
                'movie={} [watermark]; [in][watermark] overlay=main_w-overlay_w-1:1 [out]'.format(
                    self.watermark))

        # Specify the target
        args.append(target)

        return self._run_args(args)

    def to_h265(self, source: str, target: str, crf: int, preset: str, tune: str, threads: int, should_watermark=True):
        """
        Convert a video to h265

        :param source:
        :param target:

        :return:
        """

        if not os.path.exists(source):
            raise FileNotFoundError(source)

        args = [
            self.ffmpeg,

            # Overwrite if file exists
            '-y',

            # Output level (Warnings)
            '-loglevel', '24' if self.quiet else '32',

            # Input filter
            '-i', source,

            # Video
            '-c:v', 'libx264', '-crf', str(crf),

            # Audio
            '-c:a', 'libfdk_aac', '-preset', str(preset),

            # Frame rate
            '-tune', str(tune),

            # Number of threads to use
            '-threads', str(threads),
        ]

        if should_watermark:
            args.append('-vf')
            args.append(
                'movie={} [watermark]; [in][watermark] overlay=main_w-overlay_w-1:1 [out]'.format(
                    self.watermark))

        # Specify the target
        args.append(target)

        return self._run_args(args)

    def _probe_source(self, source) -> tuple:
        """
        Query the source with ffprobe returning the parse json data

        :param source:
        :return:
        """
        if not os.path.exists(source):
            raise FileNotFoundError(source)

        args = (
            self.ffprobe,

            # Display flags
            '-show_format', '-show_streams',

            # Don't output excessive information to break the json
            '-loglevel', 'quiet',

            # We want the data in json
            '-print_format', 'json',

            # file to process
            source
        )
        output = subprocess.check_output(args, env=dict(NO_COLOR='1'))

        data = json.loads(output.decode('utf-8'))

        try:
            video, audio = data['streams']
        except ValueError:
            video = data['streams'][0]
            audio = None

        return video, audio, data

    def get_encoding(self, source) -> str:
        video, audio, data = self._probe_source(source)

        return video.get('codec_name', 'unknown')

    def query_metadata(self, source):

        video, audio, data = self._probe_source(source)

        general = {
            'duration': time.strftime('%H:%M:%S', time.gmtime(int(data['format']['duration'].split('.')[0]))),
            'size_bytes': int(data['format']['size']),
            'size_mb': round(int(data['format']['size']) / (1024 * 1024), 2),
            'bit_rate': math.floor(int(data['format']['bit_rate']) / 1024)
        }

        # Don't display the video fps as a fraction
        video['fps'] = video['r_frame_rate'].split('/')[0]

        text = 'Size: {size_bytes} bytes ({size_mb} MB), duration: {duration}, avg.bitrate: {bit_rate} kb/s\n'.format(
            **general)

        if audio:
            text += 'Audio: {codec_name}, {sample_rate} Hz, {channel_layout}\n'.format(**audio)
        else:
            text += 'Audio: Not available\n'

        text += 'Video: {codec_name}, {pix_fmt}, {width}x{height}, {fps} fps'.format(
            codec_name=video.get('codec_name', 'unknown'),
            pix_fmt=video.get('pixfmt', 'unknown'),
            width=video.get('width', 'unknown'),
            height=video.get('height', 'unknown'),
            fps=video.get('fps', 'unknown'),
        )

        # Convert H:i:s to duration in seconds
        duration = sum(
            int(x) * 60 ** i for i, x in enumerate(reversed(general['duration'].split(':'))))

        return text, int(duration), general['size_bytes']

    def generate_thumbnail(self, source: str, target: str, images=25, columns=5, height=260):

        if not os.path.exists(source):
            raise FileNotFoundError(source)

        if os.path.exists(target):
            os.unlink(target)

        path = os.path.dirname(target)

        if not os.path.exists(path):
            os.makedirs(path, exist_ok=True)

        args = [
            self.vcs,

            # The file to process
            source,

            # Number of thumbnails to take
            '-n', str(images),

            # Number of columns
            '-c', str(columns),

            # Height of each shot
            '--height', str(height),

            # Target file for the jpg
            '-o', target,
        ]

        return self._run_args(args)


def transcoding_utility_factory(config) -> TranscodingUtility:
    verbose = config.FFMPEG_VERBOSE if hasattr(config, 'FFMPEG_VERBOSE') else True

    return TranscodingUtility(
        config.FFMPEG_COMMAND,
        config.FFPROBE_COMMAND,
        config.VCS_COMMAND,
        config.WATERMARK_FILE,
        quiet=not verbose
    )
