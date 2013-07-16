function HumantestController($scope,$http) {
    $scope.postfile = function () {
        $("#fileform").ajaxSubmit(function (data) {
            $scope.ShowResult = true;
            $scope.rpt = $.parseJSON(data);
            $scope.$apply();
        })
    }
}

$(function () {
    $("#fileform").ajaxForm();
})
