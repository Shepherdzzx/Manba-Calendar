package com.example.manba_alert

import bridge.Bridge
import io.flutter.embedding.android.FlutterActivity
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.plugin.common.MethodChannel

class MainActivity : FlutterActivity() {
    private val parserChannel = "com.example.manba_alert/parser"

    override fun configureFlutterEngine(flutterEngine: FlutterEngine) {
        super.configureFlutterEngine(flutterEngine)

        MethodChannel(flutterEngine.dartExecutor.binaryMessenger, parserChannel)
            .setMethodCallHandler { call, result ->
                when (call.method) {
                    "parseText" -> {
                        val input = call.argument<String>("input") ?: ""
                        val now = call.argument<String>("now") ?: ""
                        try {
                            val json = Bridge.parseText(input, now)
                            result.success(json)
                        } catch (error: Exception) {
                            result.error("PARSE_ERROR", error.message, null)
                        }
                    }

                    else -> result.notImplemented()
                }
            }
    }
}
