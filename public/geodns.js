"use strict";

(function($) {

   console.log("loading js");

   if ($('#events')) {
       var events = new EventSource('/api/events');

       // Subscribe to "log" event
       var content = $('#events');
       events.addEventListener('log',
           function(event) {
               console.log("got event", event, event.data);
               var data = JSON.parse(event.data);
               console.log("json data", data);

               if (data.level !== 'debug') {
                   var type = "info";
                   if (data.level === 'warn') {
                       type = "warning";
                   }
                   else if (data.level === 'error' || data.level === 'fatal') {
                       type = "danger";
                   }

                   $('.top-right').notify(
                       {
                           message: { text: data.lines },
                           type: type
                       }).show();
               }

               var html = '[<span class="event-' + data.level + '">' +
                   data.level + '</span>] ' +
                   data.lines + '<br>';
               content.prepend(html);
           }, false
       );
   }

})(jQuery);

var configList = ['Pops','Groups','Labels','Outages'];

var App = angular.module('geodns', ['geodnsServices', 'buttonsRadio']).
    config(['$routeProvider', function($routeProvider) {
        $routeProvider.
            when('/pops',   {templateUrl: 'views/pops.html',   controller: PopsCtrl}).
            when('/groups', {templateUrl: 'views/groups.html', controller: GroupsCtrl}).
            otherwise({redirectTo: '/pops'});
    }]);

var services = angular.module('geodnsServices', ['ngResource']);

_.each(configList, function(t) {
    services.factory(t, function($resource) {
        var api_name = t.toLowerCase();
        var r = $resource('/api/' + api_name, {});
        return r;
    });
});

angular.module('buttonsRadio', []).directive('buttonsRadio', function() {
    return {
        restrict: 'E',
        scope: { model: '=', options:'='},
        controller: function($scope){
            $scope.activate = function(option){
                $scope.model = option;
            };
        },
        template: "<button type='button' class='btn' "+
                    "ng-class='{active: option.v == model}'"+
                    "ng-repeat='option in options' "+
                    "ng-click='activate(option.v)'>{{option.l}} "+
                  "</button>"
    };
});

function PopsCtrl($scope, Pops) {
    $scope.pops = Pops.query();
    console.log("Pops", $scope.pops);
    // $scope.orderProp = 'date';
}

function GroupsCtrl($scope, Groups) {
    $scope.groups = Groups.query();
}
