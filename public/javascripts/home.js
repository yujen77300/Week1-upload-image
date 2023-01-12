const textInput = document.querySelector('#text-input')
const formFile = document.querySelector('#form-file')
const uploadButton = document.querySelector('.upload-button')
console.log(formFile)
let form = new FormData();

formFile.addEventListener('change', function (e) {
  // let form = new FormData();
  // e.target.files[0]取得使用者上傳的檔案
  form.append('form', e.target.files[0])
  console.log(form)
})

uploadButton.addEventListener('click', function () {
  console.log("近來要上傳了")
  console.log(textInput.value)
  form.append('text', textInput.value)
  uploadImage(form)
})


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
      console.log("成功了")
      console.log(result)
  
    }
  } catch (err) {
    console.log({ "error": err.message });
  }
}

