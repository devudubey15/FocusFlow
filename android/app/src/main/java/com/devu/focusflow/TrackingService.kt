package com.devu.focusflow

import android.app.*
import android.app.usage.UsageStatsManager
import android.content.Context
import android.content.Intent
import android.os.Build
import android.os.IBinder
import androidx.core.app.NotificationCompat
import java.util.*

class TrackingService : Service() {

    private val timer = Timer()
    private val INTERVAL: Long = 5000 // 5 seconds

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        createNotificationChannel()
        val notification = NotificationCompat.Builder(this, CHANNEL_ID)
            .setContentTitle("FocusFlow Tracking")
            .setContentText("FocusFlow is monitoring your focus.")
            .setSmallIcon(android.R.drawable.ic_menu_agenda)
            .build()

        startForeground(1, notification)
        startTracking()
        
        return START_STICKY
    }

    private fun startTracking() {
        timer.scheduleAtFixedRate(object : TimerTask() {
            override fun run() {
                val activeApp = getForegroundApp()
                // Sync with Supabase (simplified logic)
                println("Active App: $activeApp")
            }
        }, 0, INTERVAL)
    }

    private fun getForegroundApp(): String? {
        val usm = getSystemService(Context.USAGE_STATS_SERVICE) as UsageStatsManager
        val time = System.currentTimeMillis()
        val appList = usm.queryUsageStats(UsageStatsManager.INTERVAL_DAILY, time - 1000 * 1000, time)
        if (appList != null && appList.isNotEmpty()) {
            val sortedMap = TreeMap<Long, android.app.usage.UsageStats>()
            for (usageStats in appList) {
                sortedMap[usageStats.lastTimeUsed] = usageStats
            }
            if (!sortedMap.isEmpty()) {
                return sortedMap.get(sortedMap.lastKey())?.packageName
            }
        }
        return null
    }

    private fun createNotificationChannel() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val serviceChannel = NotificationChannel(
                CHANNEL_ID,
                "FocusFlow Tracking Channel",
                NotificationManager.IMPORTANCE_DEFAULT
            )
            val manager = getSystemService(NotificationManager::class.java)
            manager.createNotificationChannel(serviceChannel)
        }
    }

    override fun onBind(intent: Intent?): IBinder? = null

    companion object {
        const val CHANNEL_ID = "FocusFlowServiceChannel"
    }
}
