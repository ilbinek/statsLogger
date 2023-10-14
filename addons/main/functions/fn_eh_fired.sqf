params ["_unit", "_weapon", "_muzzle", "_mode", "_ammo", "_magazine", "_projectile", "_gunner"];

private _tmp = [getPlayerUID _unit];
"stats_logger" callExtension [":SHOT:", _tmp];
