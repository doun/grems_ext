function HumantestController($scope,$http) {
    $scope.post_reply = function(){

    }
    $scope.postfile= function(){
        data = $("form").serialize()
        $http.post("http://localhost:909/parser",data).success(function (data) {
            alert(data)
        }).error(function (data) {
            alert("failed" + data)
        })
    }
}
