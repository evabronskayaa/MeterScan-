const token = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTg5MTc2MzYsIm9yaWdfaWF0IjoxNjk4OTE0MDM2LCJ1c2VyX2lkIjoxfQ.LYB40lhL6RSM6SEyST2tVGopwSgXjsPbk8ZhFjK7j8M"

//export let number = "228"
export let number

export const uploadImage = async (image) => {
    // console.error(image instanceof Blob)
    // let form = new FormData()
    // form.append('files', image)

    // const result = await fetch("http://localhost:8080/api/v1/predict", {
    //     method: 'POST',
    //     body: form,
    //     headers: {
    //         'Authorization': token
    //     }
    // })

    // if (result.status === 200)
    // number = result.body
    // else
    number = ["25", "27"]
    console.log(typeof number)
    
    //console.log(result)
}
