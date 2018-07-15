package com.example.cryptoalert

import okhttp3.RequestBody
import retrofit2.Call
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.POST

interface ApiService {

  @POST("/btc-pref")
  fun saveBTCLimit(@Body body: RequestBody): Call<String>

  @POST("/eth-pref")
  fun saveETHLimit(@Body body: RequestBody): Call<String>

  @GET("/fetch-values")
  fun getValues():Call<String>

}
