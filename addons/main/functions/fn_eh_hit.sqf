params ["_unit", "_causedBy", "_damage", "_instigator"]; 
private _tmp = [getPlayerUID _instigator]; 
"stats_logger" callExtension [":HIT:", _tmp]; 
