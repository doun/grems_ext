function HumantestController($scope,$http) {
}

$(function () {
    $("#fileform").ajaxForm()
    $("#fileform").on('submit',function () {
        $(this).ajaxSubmit(function (data) {
            alert(data)
        })
        return false
    })
})
