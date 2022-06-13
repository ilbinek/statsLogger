/**
 * @file main.cpp
 * @author Sotirios Pupakis (https://github.com/ilbinek)
 * @brief Main file of the Stats Plugin extension taking care of everything
 * @version 0.1
 * @date 2022-06-13
 * 
 */

#include <algorithm>
#include <cstring>
#include <ctime>
#include <fstream>
#include <iostream>
#include <iterator>
#include <sstream>
#include <string>
#include <vector>
#include <iomanip>

#include "kill.hpp"
#include "mission.hpp"
#include "player.hpp"

#define PROTECTION

extern "C" {
__attribute__((dllexport)) void RVExtension(char *output, int outputSize, const char *function);
__attribute__((dllexport)) int RVExtensionArgs(char *output, int outputSize, const char *function, const char **argv, int argc);
__attribute__((dllexport)) void RVExtensionVersion(char *output, int outputSize);
}

void split(const std::string &str, const std::string &delim, std::vector<std::string> &parts);
void writeMission(std::ofstream& outfile);
void writePlayers(std::ofstream& outfile);
void writeKills(std::ofstream& outfile);
//void writeObjectives(std::ofstream& outfile);
std::string escape_json(const std::string &s);
std::string getDate();

std::vector<Player> players;
std::vector<Kill> kills;
Mission mission;

void RVExtension(char *output, int outputSize, const char *function) {
    try {
        std::strncpy(output, function, outputSize - 1);
        char* errmsg = "WRONG ARGUMENT NUMBER";

        // Parse C-string into C++ string
        std::string str = function;
        std::string delim = "::";
        std::vector<std::string> vec;
        // Parse them into vector
        split(str, delim, vec);
        
        if (vec.at(0) == "RESET") {
            // Reset everthing and forget it. I think it should free it in memory, not sure though
            mission.clear();
            players.clear();
            kills.clear();
        } else if (vec.at(0) == "MISSION") {
            #ifdef PROTECTION
            if (vec.size() != 6) {
                std::strncpy(output, errmsg, strlen(errmsg) - 1);
                return;
            }
            #endif
            time_t now = time(0);
            char *dt = ctime(&now);
            mission.missionSet(vec.at(1), vec.at(2), vec.at(3), vec.at(4), vec.at(5), dt);
        } else if (vec.at(0) == "WIN") {
            #ifdef PROTECTION
            if (vec.size() != 5) {
                std::strncpy(output, errmsg, strlen(errmsg) - 1);
                return;
            }
            #endif
            mission.setVictory(vec.at(1));
            mission.setEndTime(vec.at(2));
            mission.setScoreBlue(vec.at(3));
            mission.setScoreRed(vec.at(4));
        } else if (vec.at(0) == "PLAYER") {
            #ifdef PROTECTION
            if (vec.size() != 7) {
                std::strncpy(output, errmsg, strlen(errmsg) - 1);
                return;
            }
            #endif
            // Iterate through all players and find out if he is already logged
            std::string uid = vec.at(1);
            bool found = false;
            for (int i = 0; i < players.size(); i++) {
                if (players.at(i).getUID() == uid) {
                    found = true;
                    players.at(i).setRole(vec.at(3));
                    players.at(i).setClassName(vec.at(4));
                    players.at(i).setSide(vec.at(5));
                    players.at(i).setSquad(vec.at(6));
                    break;
                }
            }
            // Return if found
            if (found) {
                return;
            }
            // Create new object
            players.push_back(Player(vec.at(1), vec.at(2), vec.at(3), vec.at(4), vec.at(5), vec.at(6)));
        } else if (vec.at(0) == "KILL") {
            #ifdef PROTECTION
            if (vec.size() != 6) {
                std::strncpy(output, errmsg, strlen(errmsg) - 1);
                return;
            }
            #endif
            kills.push_back(Kill(vec.at(1), vec.at(2), vec.at(3), vec.at(4), vec.at(5)));
        } else if (vec.at(0) == "SHOT") {
            #ifdef PROTECTION
            if (vec.size() != 2) {
                std::strncpy(output, errmsg, strlen(errmsg) - 1);
                return;
            }
            #endif
            std::string uid = vec.at(1);

            for (int i = 0; i < players.size(); i++) {
                if (players.at(i).getUID() == uid) {
                    players.at(i).shot();
                    break;
                }
            }
        } else if (vec.at(0) == "HIT") {
            #ifdef PROTECTION
            if (vec.size() != 2) {
                std::strncpy(output, errmsg, strlen(errmsg) - 1);
                return;
            }
            #endif
            std::string uid = vec.at(1);

            for (int i = 0; i < players.size(); i++) {
                if (players.at(i).getUID() == uid) {
                    players.at(i).hit();
                    break;
                }
            }
        } else if (vec.at(0) == "EXPORT") {
            #ifdef PROTECTION
            if (vec.size() != 4) {
                std::strncpy(output, errmsg, strlen(errmsg) - 1);
                return;
            }
            #endif
            if (vec.at(1) == "REALLY" && vec.at(2) == "REALLY" && vec.at(3) == "REALLY") {
                // Get current date to make it the file name
                std::string fname = getDate();
                fname = fname + "_" + mission.getName() + ".json";
                // Export into json
                std::ofstream outfile ("/stats-output/" + fname);
                // Start of json
                outfile << "{\n";
                // Write mission
                writeMission(outfile);
                // Write Playres
                writePlayers(outfile);
                // Write kills
                writeKills(outfile);
                // End of json
                outfile << "}\n";
            }
        }
    } catch (...) {
        // Something happened, it's terrible, but we can't throw an error as the Arma server would die too.
        // TODO add a log file and log this event
    }
}

int RVExtensionArgs(char *output, int outputSize, const char *function, const char **argv, int argc) {
    std::stringstream sstream;
    for (int i = 0; i < argc; i++) {
        sstream << argv[i];
    }
    std::strncpy(output, sstream.str().c_str(), outputSize - 1);
    return 0;
}

void RVExtensionVersion(char *output, int outputSize) {
    std::strncpy(output, "Stats Extension - Version 0.1", outputSize - 1);
}

// taken from https://stackoverflow.com/a/289365
void split(const std::string &str, const std::string &delim, std::vector<std::string> &parts) {
    size_t last = 0;
    size_t next = 0;
    while ((next = str.find(delim, last)) != std::string::npos) {
        parts.push_back(escape_json(str.substr(last, next - last)));
        last = next + delim.length();
    }
    parts.push_back(escape_json(str.substr(last)));
}

void writeMission(std::ofstream& outfile) {
    // Mission related things
    outfile << "\"worldname\":\"" << mission.getMap() << "\",\n";
    outfile << "\"missionName\":\"" << mission.getName() << "\",\n";
    outfile << "\"missionAuthor\":\"" << mission.getAuthor() << "\",\n";
    // TODO?
    outfile << "\"missionType\":\"public\"" << ",\n";
    outfile << "\"victory\":\"" << mission.getVictory() << "\",\n";
    outfile << "\"missionStart\":\"" << mission.getMissionStart() << "\",\n";
    outfile << "\"missionEnd\":\"" << mission.getMissionEnd() << "\",\n";
}

void writePlayers(std::ofstream& outfile) {
    // All player
    // Start of an arary
    outfile << "\"players\":[" << "\n";
    bool first = true;
    // Data for each player
    for (auto player : players) {
        if (!first) {
            outfile << ",\n";
        }
        first = false;
        outfile << "\t{\"uid\":\"" << player.getUID() << "\", ";
        outfile << "\"name\":\"" << player.getName() << "\", ";
        outfile << "\"side\":\"" << player.getSide() << "\", ";
        outfile << "\"shots\":\"" << player.getShots() << "\", ";
        outfile << "\"hits\":\"" << player.getHits() << "\", ";
        outfile << "\"squad\":\"" << player.getSquad() << "\", ";
        outfile << "\"role\":\"" << player.getRole() << "\",";
        outfile << "\"class\":\"" << player.getClassName() << "\"}";
    }
    outfile << "\n";
    // End of an arary
    outfile << "]," << "\n";
}

void writeKills(std::ofstream& outfile) {
    // All kills
    // Start of an arary
    outfile << "\"kills\":[" << "\n";
    bool first = true;
    // Data for each kill
    for (auto kill : kills) {
         if (!first) {
            outfile << ",\n";
        }
        first = false;
        outfile << "\t{\"time\":\"" << kill.getTime() << "\", ";
        outfile << "\"victim\":\"" << kill.getKilledUID() << "\", ";
        outfile << "\"killer\":\"" << kill.getKillerUID() << "\", ";
        outfile << "\"weapon\":\"" << kill.getWeapon() << "\", ";
        outfile << "\"distance\":\"" << kill.getDistance() << "\"}";
    }
    outfile << "\n";
    // End of an arary
    outfile << "]" << "\n";
}

std::string getDate() {
    std::time_t t = std::time(0);
    std::tm* now = std::localtime(&t);
    std::string month = (now->tm_mon + 1 < 10 ) ? "0" + std::to_string(now->tm_mon + 1) : std::to_string(now->tm_mon + 1);
    return std::to_string(now->tm_year + 1900) + "-" + month + "-" + std::to_string(now->tm_mday) + "-" + std::to_string(now->tm_hour) + "-" + std::to_string(now->tm_min);
}

std::string escape_json(const std::string &s) {
    std::ostringstream o;
    for (auto c = s.cbegin(); c != s.cend(); c++) {
        switch (*c) {
        case '"': o << "\\\""; break;
        case '\\': o << "\\\\"; break;
        case '\b': o << "\\b"; break;
        case '\f': o << "\\f"; break;
        case '\n': o << "\\n"; break;
        case '\r': o << "\\r"; break;
        case '\t': o << "\\t"; break;
        default:
            if ('\x00' <= *c && *c <= '\x1f') {
                o << "\\u"
                  << std::hex << std::setw(4) << std::setfill('0') << static_cast<int>(*c);
            } else {
                o << *c;
            }
        }
    }
    return o.str();
}
