#include "mission.hpp"

void Mission::missionSet(std::string name, std::string map, std::string author, std::string missionType, std::string missionStart, std::string date) {
    this->name = name;
    this->map = map;
    this->author = author;
    this->missionType = missionType;
    this->missionStart = missionStart;
    this->date = date;
}

void Mission::clear() {
    this->name = "";
    this->map = "";
    this->author = "";
    this->missionType = "";
    this->missionStart = "";
    this->date = "";
    this->victory = "";
    this->missionEnd = "";
}
