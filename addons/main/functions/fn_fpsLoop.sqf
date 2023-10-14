diag_log(text ('[STATS] FPS LOOP CALLED'));
// Loop to output server FPS every second
[] spawn {
	diag_log('[STATS] FPS OOP Spawned');
	while {true} do {
		private _tmp = [diag_fps];
		"stats_logger" callExtension [":FPS:", _tmp];
		sleep 1;
	};
};

diag_log(text ('[STATS] FPS logger started '));
