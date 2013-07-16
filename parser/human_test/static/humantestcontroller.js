function HumantestController($scope,$http) {
    $scope.postfile = function () {
        $("#fileform").ajaxSubmit(function (data) {
            alert(data)
        })
    }
}

$(function () {
    $("#fileform").ajaxForm()
})
