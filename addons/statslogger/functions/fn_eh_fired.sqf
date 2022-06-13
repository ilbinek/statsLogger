params ["_unit", "_weapon", "_muzzle", "_mode", "_ammo", "_magazine", "_projectile", "_gunner"];
_tmp = ["SHOT", getPlayerUID _unit] joinString "::";
diag_log(text ('[STATS] ' + _tmp));

"Stats" callExtension _tmp;
