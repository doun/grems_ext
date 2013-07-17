function HumantestController($scope,$http) {
    $scope.HasResult = false;
    $scope.postfile = function () {
        $("#fileform").ajaxSubmit({
            success:function (data) {
                $scope.$apply(function () {
                    rpt = $.parseJSON(data);
                    if (rpt.error == "true") {
                        alert("文件解析出错:"+rpt.description);
                    }else{
                        $scope.rpt = rpt;
                        $scope.HasResult = true;
                    }
                });
                $("#rsp_wrapper").draggable({
                    stop:function (event,ui) {
                        $(this).css("position","fixed");
                    }
                })
            },
            dataType:'text'})
    }
}

$(function () {
    $("#fileform").ajaxForm();
})
