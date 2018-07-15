package com.example.cryptoalert

import android.content.Context
import android.content.SharedPreferences

class Prefs (context: Context) {
    val PREFS_FILENAME = "com.example.coinalert.prefs"
    val DEVICE_UUID = "device_uuid"
    val prefs: SharedPreferences = context.getSharedPreferences(PREFS_FILENAME, 0);

    var deviceUuid: String
        get() = prefs.getString(DEVICE_UUID, "")
        set(value) = prefs.edit().putString(DEVICE_UUID, value).apply()

}