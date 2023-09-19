//params ["_unit", "_source", "_damage", "_instigator"];
private _tmp = [getPlayerUID (_this select 1)];
diag_log(text ('[STATS] ' + str(_tmp)));
"stats_logger" callExtension [":HIT:", _tmp];
