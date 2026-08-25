export function buildPieceworkQrPayload({ workOrderId, processId, stationId }) {
  return {
    type: 'piecework',
    version: 1,
    workOrderId,
    processId,
    stationId,
  }
}

export function scanQrCode() {
  return new Promise((resolve) => {
    // H5 端 uni.scanCode 不可用，先返回明确提示，方便页面流程预览。
    // #ifdef H5
    resolve({
      text: 'H5 暂不调用摄像头扫码，请直接输入订单号预览流程，真机小程序/App 使用扫码。',
      raw: null,
    })
    // #endif

    // #ifndef H5
    uni.scanCode({
      onlyFromCamera: true,
      success: (res) => {
        resolve({
          text: res.result,
          raw: res,
        })
      },
      fail: (error) => {
        resolve({
          text: error.errMsg || '扫码失败',
          raw: error,
        })
      },
    })
    // #endif
  })
}
