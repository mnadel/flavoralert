Array.prototype.getUnique = function () {
   var u = {}, a = [];
   for(var i = 0, l = this.length; i < l; ++i){
      if(u.hasOwnProperty(this[i])) {
         continue;
      }
      a.push(this[i]);
      u[this[i]] = 1;
   }
   return a;
}

var app = angular.module("flavoralert", []);

app.controller("Alerting", function ($scope, $http, $location) {
    $http({
        method: "GET",
        url: $location.protocol() + "://" + $location.host() + ":" + $location.port() + "/version/0/meta"
    }).success(function (resp) {
        $scope.meta = resp.data;
    });

    $scope.save = function () {
        alert("Saving!");
    }
});

app.controller("CurrentFlavors", function ($scope, $http, $location) {
    $http({
        method: "GET",
        url: $location.protocol() + "://" + $location.host() + ":" + $location.port() + "/version/0/current"
    }).success(function (resp) {
        $scope.locs = resp.data;
        $scope.flavors = {};

        Object.getOwnPropertyNames(resp.data).forEach(function (loc) {
            resp.data[loc].forEach(function (flavor) {
                if (!$scope.flavors[flavor]) {
                    $scope.flavors[flavor] = [loc];
                } else {
                    $scope.flavors[flavor].push(loc);
                }
            });
        });
    });
});

app.controller("AllFlavors", function ($scope, $http, $location) {
    $http({
        method: "GET",
        url: $location.protocol() + "://" + $location.host() + ":" + $location.port() + "/version/0/all"
    }).success(function (resp) {
        $scope.flavors = resp.data.getUnique();
    });
});