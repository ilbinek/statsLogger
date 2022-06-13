#ifndef MISSION_H
#define MISSION_H

#include <string>

class Mission {
    private:
        std::string name;
        std::string map;
        std::string author;
        std::string missionType;
        std::string victory;
        std::string missionStart;
        std::string missionEnd;
        std::string date;
        std::string scoreBlue;
        std::string scoreRed;

    public:
        void missionSet(std::string name, std::string map, std::string author, std::string missionType, std::string missionStart, std::string date);
        void setVictory(std::string victory) {this->victory = victory;};
        void setEndTime(std::string endTime) {missionEnd = endTime;};
        void setScoreBlue(std::string scoreBlue) {this->scoreBlue = scoreBlue;};
        void setScoreRed(std::string scoreRed) {this->scoreRed = scoreRed;};
        void clear();
        std::string getName() {return name;};
        std::string getMap() {return map;};
        std::string getAuthor() {return author;};
        std::string getMissionType() {return missionType;};
        std::string getVictory() {return victory;};
        std::string getMissionStart() {return missionStart;};
        std::string getMissionEnd() {return missionEnd;};
};

#endif
