const textInput = document.querySelector('#text-input')
const formFile = document.querySelector('#form-file')
const uploadButton = document.querySelector('.upload-button')
const container = document.querySelector('.container')
const row = document.querySelector('.row')
console.log(formFile)
let form = new FormData();
renderHomePage()

formFile.addEventListener('change', function (e) {
  // let form = new FormData();
  // e.target.files[0]取得使用者上傳的檔案
  form.append('form', e.target.files[0])
  // console.log(form)
})

uploadButton.addEventListener('click', function () {
  // console.log("近來要上傳了")
  // console.log(textInput.value)
  form.append('text', textInput.value)
  uploadImage(form)
  form = new FormData();
})


function renderHomePage() {
  fetch(
    '/api/allfile'
  ).then((response) => {
    return response.json();
  }).then((data) => {
    // console.log("進來選染一開始的資料囉~")
    // console.log(data)
    // console.log(data[0]["ImageUrl"])
    // console.log(data[1]["Text"])
    // console.log(data[2]["Id"])
    data.forEach((element) => {
      const resultSection = document.createElement('div')
      resultSection.className = "result-section"
      const hr = document.createElement('hr')
      resultSection.appendChild(hr)
      const resultInfo = document.createElement('div')
      resultInfo.className = "result"
      resultSection.appendChild(resultInfo)
      const resultText = document.createElement('div')
      resultText.className = "result-text"
      resultText.textContent = element["Text"]
      resultInfo.appendChild(resultText)
      const resultImage = document.createElement('div')
      resultImage.className = "result-img"
      resultImage.style.cssText = `background-image: url(${element["ImageUrl"]});height:150px;width:200px; background-size:cover;background-position:center center;`
      resultInfo.appendChild(resultImage)
      row.insertAdjacentElement('afterend', resultSection)
    });
  })
}


async function uploadImage(form) {
  let url = "/api/upload/image"
  let options = {
    body: form,
    method: "POST",
  }
  try {
    let response = await fetch(url, options);
    let result = await response.json();
    if (response.status === 200) {
      // console.log("成功了")
      // console.log(response)
      // console.log("結果")
      // console.log(result)
      const resultSection = document.createElement('div')
      resultSection.className = "result-section"
      const hr = document.createElement('hr')
      resultSection.appendChild(hr)
      const resultInfo = document.createElement('div')
      resultInfo.className = "result"
      resultSection.appendChild(resultInfo)
      const resultText = document.createElement('div')
      resultText.className = "result-text"
      resultText.textContent = result["Text"]
      resultInfo.appendChild(resultText)
      const resultImage = document.createElement('div')
      resultImage.className = "result-img"
      resultImage.style.cssText = `background-image: url(${result["CloudFrontUrl"]});height:150px;width:200px; background-size:cover;background-position:center center;`
      resultInfo.appendChild(resultImage)
      row.insertAdjacentElement('afterend', resultSection)
      textInput.value = ""
      formFile.value = ""
    }
  } catch (err) {
    console.log({ "error": err.message });
  }
}
