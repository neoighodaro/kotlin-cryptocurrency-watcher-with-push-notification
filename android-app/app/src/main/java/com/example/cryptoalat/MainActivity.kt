package com.example.cryptoalat

import android.os.Bundle
import okhttp3.MediaType
import okhttp3.RequestBody
import org.json.JSONObject
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response
import android.support.v7.app.AlertDialog
import android.support.v7.app.AppCompatActivity
import android.util.Log
import android.view.LayoutInflater
import android.widget.Button
import android.widget.EditText
import com.pusher.pushnotifications.PushNotifications
import kotlinx.android.synthetic.main.activity_main.*
import okhttp3.OkHttpClient
import retrofit2.Retrofit
import retrofit2.converter.scalars.ScalarsConverterFactory

class MainActivity : AppCompatActivity() {

  private val retrofit: ApiService by lazy {
    val httpClient = OkHttpClient.Builder()
    val builder = Retrofit.Builder()
        .baseUrl("http://10.0.2.2:4000/")
        .addConverterFactory(ScalarsConverterFactory.create())

    val retrofit = builder
        .client(httpClient.build())
        .build()
    retrofit.create(ApiService::class.java)
  }

  override fun onCreate(savedInstanceState: Bundle?) {
    super.onCreate(savedInstanceState)
    setContentView(R.layout.activity_main)
    fetchCurrentPrice()
    setupPushNotifications()
    setupClickListeners()
  }


  private fun fetchCurrentPrice() {
    retrofit.getValues().enqueue(object: Callback<String> {
      override fun onResponse(call: Call<String>?, response: Response<String>?) {
        val jsonObject = JSONObject(response!!.body())
        bitcoinValue.text = "1 BTC = $"+ jsonObject.getJSONObject("BTC").getString("USD")
        etherumValue.text = "1 ETH = $"+ jsonObject.getJSONObject("ETH").getString("USD")
      }

      override fun onFailure(call: Call<String>?, t: Throwable?) {
        Log.e("MainActivity",t!!.localizedMessage)
      }
    })
  }

  private fun setupPushNotifications() {
    PushNotifications.start(applicationContext, "PUSHER_BEAMS_INSTANCE_ID")
    PushNotifications.subscribe("crypto")
  }


  private fun setupClickListeners() {
    bitcoinValue.setOnClickListener {
      createDialog("BTC")
    }
    etherumValue.setOnClickListener {
      createDialog("ETH")
    }

  }


  private fun createDialog(source:String){

    val builder: AlertDialog.Builder = AlertDialog.Builder(this)
    val view = LayoutInflater.from(this).inflate(R.layout.alert_layout,null);
    builder.setTitle("Set limits")
        .setMessage("Notifications will be sent to you when the value exceeds or goes below the maximum and minimum values")
        .setView(view)

    val minEditText:EditText = view.findViewById(R.id.minimumValue)
    val maxEditText:EditText = view.findViewById(R.id.maximumValue)
    val dialog = builder.create()

    view.findViewById<Button>(R.id.save).setOnClickListener {
      if (source=="BTC"){
        saveBTCPref(minEditText.text.toString(),maxEditText.text.toString())
      } else {
        saveETHPref(minEditText.text.toString(),maxEditText.text.toString())
      }
      dialog.dismiss()
    }
    dialog.show()
  }


  private fun saveBTCPref(min:String, max:String){
    val jsonObject = JSONObject()
    jsonObject.put("minBTC",min)
    jsonObject.put("maxBTC",max)

    val body = RequestBody.create(
        MediaType.parse("application/json"),
        jsonObject.toString()
    )

    retrofit.saveBTCLimit(body).enqueue(object: Callback<String> {
      override fun onResponse(call: Call<String>?, response: Response<String>?) {}
      override fun onFailure(call: Call<String>?, t: Throwable?) {}
    })
  }

  private fun saveETHPref(min:String, max:String){
    val jsonObject = JSONObject()
    jsonObject.put("minETH",min)
    jsonObject.put("maxETH",max)

    val body = RequestBody.create(
        MediaType.parse("application/json"),
        jsonObject.toString()
    )

    retrofit.saveETHLimit(body).enqueue(object: Callback<String> {
      override fun onResponse(call: Call<String>?, response: Response<String>?) {}
      override fun onFailure(call: Call<String>?, t: Throwable?) {}
    })
  }


}
