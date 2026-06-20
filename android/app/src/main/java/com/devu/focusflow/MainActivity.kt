package com.devu.focusflow

import android.content.Intent
import android.os.Bundle
import android.provider.Settings
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            MaterialTheme {
                Surface(
                    modifier = Modifier.fillMaxSize(),
                    color = MaterialTheme.colorScheme.background
                ) {
                    FocusFlowScreen()
                }
            }
        }
    }

    @Composable
    fun FocusFlowScreen() {
        var isTracking by remember { mutableStateOf(false) }

        Column(
            modifier = Modifier.fillMaxSize().padding(24.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.Center
        ) {
            Text(
                text = "FocusFlow",
                style = MaterialTheme.typography.headlineLarge,
                modifier = Modifier.padding(bottom = 8.dp)
            )
            Text(
                text = "Unified Cross-Device Time Tracker",
                style = MaterialTheme.typography.bodyMedium,
                color = MaterialTheme.colorScheme.secondary,
                modifier = Modifier.padding(bottom = 48.dp)
            )

            Button(
                onClick = {
                    if (!isTracking) {
                        startService(Intent(this@MainActivity, TrackingService::class.java))
                    } else {
                        stopService(Intent(this@MainActivity, TrackingService::class.java))
                    }
                    isTracking = !isTracking
                },
                modifier = Modifier.fillMaxWidth().height(56.dp)
            ) {
                Text(if (isTracking) "Stop Tracking" else "Start Tracking")
            }

            Spacer(modifier = Modifier.height(16.dp))

            OutlinedButton(
                onClick = {
                    startActivity(Intent(Settings.ACTION_USAGE_ACCESS_SETTINGS))
                },
                modifier = Modifier.fillMaxWidth().height(56.dp)
            ) {
                Text("Grant Usage Access")
            }
        }
    }
}
