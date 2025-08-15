import 'dart:math';

import 'package:flutter/material.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class VolumeControl extends StatefulWidget {
  const VolumeControl({super.key});

  @override
  State<VolumeControl> createState() => _VolumeControlState();
}

class _VolumeControlState extends State<VolumeControl> {
  final WebSocketChannel _volumeChannel = WebSocketChannel.connect(
    Uri.parse('ws://localhost:8080/ws'),
  );

  void _incrementVolume() {
    if (_volumeChannel.closeCode != null) return;

    _volumeChannel.sink.add('volume_up');
  }

  void _decrementVolume() {
    if (_volumeChannel.closeCode != null) return;

    _volumeChannel.sink.add('volume_down');
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
        VolumeLevelSlider(volumeLevel: _volumeChannel.stream),
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          spacing: 20.0,
          children: <Widget>[
            ElevatedButton.icon(
              onPressed: _decrementVolume,
              icon: const Icon(Icons.remove),
              label: const Text('Decrement'),
            ),
            ElevatedButton.icon(
              onPressed: _incrementVolume,
              icon: const Icon(Icons.add),
              label: const Text('Increment'),
            ),
          ],
        ),
      ],
    );
  }
}

class VolumeLevelSlider extends StatelessWidget {
  const VolumeLevelSlider({required this.volumeLevel, super.key});

  final Stream<dynamic> volumeLevel;

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
            Text('$level%'),
            Slider.adaptive(
              value: level.toDouble(),
              min: 0,
              max: 100,
              onChanged: null,
            ),
          ],
        );
      },
    );
  }
}
