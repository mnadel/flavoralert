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

function CreateUrl(location, paths) {
    if (!(paths instanceof Array)) {
        paths = [paths];
    }

    return location.protocol() + "://" + location.host() + ":" + location.port() + "/" + paths.join("/");
}

var app = angular.module("flavoralert", []);

app.controller("Alerting", function ($scope, $http, $location) {
    $http({
        method: "GET",
        url: CreateUrl($location, "version/0/meta")
    }).success(function (resp) {
        $scope.meta = resp.data;
        $scope.alerts = {};
    });

    $scope.save = function (val) {
        Object.getOwnPropertyNames($scope.alerts).forEach(function (flavor) {
            var newState = $scope.alerts[flavor];
            var action = newState ? "add" : "remove";

            $http({
                method: "POST",
                url: CreateUrl($location, ["version/0/alert", action, flavor])
            }).success(function (resp) {
                $scope.meta = resp.data;
                $scope.alerts = {};
            });
        });
    };
});

app.controller("CurrentFlavors", function ($scope, $http, $location) {
    $http({
        method: "GET",
        url: CreateUrl($location, "version/0/current")
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
        url: CreateUrl($location, "version/0/all")
    }).success(function (resp) {
        var validFlavors = function (i) {
            return i != "Unknown" && i != "Coming Soon";
        };

        $scope.flavors = resp.data.filter(validFlavors).getUnique();
    });
});