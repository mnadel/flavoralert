var app = angular.module("flavoralert", []);

app.controller("Flavors", function ($scope, $http, $location) {
    $http({
        method: "GET",
        url: $location.protocol() + "://" + $location.host() + ":" + $location.port() + "/version/0/list"
    }).success(function (resp) {
        $scope.flavors = resp.data;
    });
});