//params ["_unit", "_source", "_damage", "_instigator"];
_tmp = ["HIT", getPlayerUID (_this select 1)] joinString "::";
diag_log(text ('[STATS] ' + _tmp));
"Stats" callExtension _tmp;
