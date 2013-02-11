"use strict";

(function($) {

    console.log("loading js");

    if ($('#events')) {
        var events = new EventSource('/api/events');

        // Subscribe to "log" event
        var content = $('#events');
        events.addEventListener('log',
            function(event) {
                // console.log("got event", event, event.data);
                var data = JSON.parse(event.data);
                // console.log("json data", data);

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

var configList = ['Nodes','Groups','Labels','Monitor','MonitorStatus'];

var controllers = {
    "Nodes":   NodesCtrl,
    "Groups":  GroupsCtrl,
    "Labels":  LabelsCtrl,
    "Monitor": MonitorCtrl
};

var App = angular.module('geodns', ['geodnsServices', 'buttonsRadio']).
    config(['$routeProvider', function($routeProvider) {
        _.each(configList, function(t) {
            if (!controllers[t]) { return; }
            var api_name = t.toLowerCase();
            $routeProvider.when('/' + api_name,
                {
                    templateUrl: 'views/' + api_name + '.html',
                    controller: controllers[t]
                }
            )
        });
        $routeProvider.otherwise({redirectTo: '/nodes'});
    }]);

var services = angular.module('geodnsServices', ['ngResource']);

_.each(configList, function(t) {
    if (t === 'MonitorStatus') {
        services.factory(t, function($resource) {
            return $resource('/api/monitor/:monitor/status',
                {monitor: '@monitor'}, {
            });
        });
    }
    else {
        services.factory(t, function($resource) {
            var api_name = t.toLowerCase();
            var r = $resource('/api/' + api_name, {});
            return r;
        });
    }
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

function NodesCtrl($scope, Nodes) {
    $scope.nodes = Nodes.query();
    console.log("Nodes", $scope.nodes);
    // $scope.orderProp = 'date';
}

function GroupsCtrl($scope, Groups) {
    $scope.groups = Groups.query();
}

function MonitorCtrl($scope, Monitor, MonitorStatus) {
    console.log("monitors");
    $scope.monitor_list = Monitor.get({}, function(m) {
        console.log("got monitor", m);
        $scope.state = {};
        _.each(m.monitor, function(mon) {
            var mon = m.monitor[0];
            console.log("getting status for", mon);
            $scope.state[mon] = MonitorStatus.get({ monitor: mon },
                function(state) {
                    console.log("got state", state);
                }
            );
            // $scope.xstate = $scope.state[m];
        });
    });

    // console.log("setup watcher");
   /* $scope.$watch("list",
    function(value) {
        console.log("watching", value);
        if (value && value.length > 0) {
            console.log("setting monitor1");
            $scope.monitor1 = $scope.list.monitors[0];
            $scope.state = Monitors.query();
        }
    }, true);

    setTimeout(function() {console.log("MONINT", $scope.monitor);}, 500);
    */
    // $scope.outages = Outages.query();
}

function LabelsCtrl($scope, Labels) {
    $scope.labels = Labels.query();
}
