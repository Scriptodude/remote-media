import 'dart:math';

import 'package:flutter/material.dart';
import 'package:flutter/foundation.dart' as Foundation;
import 'package:web_socket_channel/web_socket_channel.dart';

class AudioControl extends StatefulWidget {
  const AudioControl({super.key});

  @override
  State<AudioControl> createState() => _AudioControlState();
}

class _AudioControlState extends State<AudioControl> {
  final WebSocketChannel _volumeChannel = WebSocketChannel.connect(
    Uri.parse(
      'ws://${Uri.base.host}:${Foundation.kDebugMode ? 8080 : Uri.base.port}/ws',
    ),
  );

  void _incrementVolume() => _sendMessage("volume_up");
  void _decrementVolume() => _sendMessage("volume_down");
  void _playNext() => _sendMessage("play_next");
  void _playPrevious() => _sendMessage("play_previous");

  void _sendMessage(String message) {
    if (_volumeChannel.closeCode != null) return;

    _volumeChannel.sink.add(message);
  }

  @override
  void dispose() {
    if (_volumeChannel.closeCode == null) {
      _volumeChannel.sink.close(1000, 'App closed');
    }
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      spacing: 10.0,
      children: <Widget>[
        VolumeLevelSlider(
          volumeLevel: _volumeChannel.stream,
          setVolumeFn: (value) => _sendMessage("set_volume=${value.toInt()}"),
        ),
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          spacing: 20.0,
          children: <Widget>[
            IconButton(
              onPressed: _playPrevious,
              icon: const Icon(Icons.skip_previous),
            ),
            IconButton(
              onPressed: _decrementVolume,
              icon: const Icon(Icons.remove),
            ),
            IconButton(
              onPressed: _incrementVolume,
              icon: const Icon(Icons.add),
            ),
            IconButton(onPressed: _playNext, icon: const Icon(Icons.skip_next)),
          ],
        ),
      ],
    );
  }
}

class VolumeLevelSlider extends StatelessWidget {
  const VolumeLevelSlider({
    required this.volumeLevel,
    this.setVolumeFn,
    super.key,
  });

  final Stream<dynamic> volumeLevel;
  final void Function(double value)? setVolumeFn;

  @override
  Widget build(BuildContext context) {
    return StreamBuilder<dynamic>(
      stream: volumeLevel,
      builder: (BuildContext context, AsyncSnapshot<dynamic> snapshot) {
        if (snapshot.hasError) {
          return Text(snapshot.error.toString());
        }

        int level = 0;

        if (snapshot.hasData) {
          if (snapshot.data is String && (snapshot.data as String).isNotEmpty) {
            try {
              level = int.parse(snapshot.data as String);
            } catch (e) {
              print(e);
            }
          } else if (snapshot.data is int) {
            level = snapshot.data as int;
          }
        }

        level = min(max(level, 0), 100);

        return Column(
          children: [
            Text('volume: $level%'),
            Slider.adaptive(
              allowedInteraction: SliderInteraction.tapAndSlide,
              value: level.toDouble(),
              min: 0,
              max: 100,
              onChanged: setVolumeFn,
            ),
          ],
        );
      },
    );
  }
}
