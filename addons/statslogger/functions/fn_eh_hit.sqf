params ["_unit", "_causedBy", "_damage", "_instigator"];
diag_log("HIT HIT HIT HIT");
private _tmp = [getPlayerUID _instigator];
"stats_logger" callExtension [":HIT:", _tmp];
