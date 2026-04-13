#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
CREATE_URL="${BASE_URL}/api/movie/create"
TOKEN="${TOKEN:-}"

# If TOKEN is set, requests include Bearer auth; otherwise they are sent without Authorization.
auth_header=()
if [ -n "$TOKEN" ]; then
  auth_header=(-H "Authorization: Bearer ${TOKEN}")
fi

post_movie() {
  local name="$1"
  local genre="$2"
  local description="$3"
  local ratings="$4"

  curl -s -X POST "$CREATE_URL" \
    "${auth_header[@]}" \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"${name}\",\"genre\":\"${genre}\",\"description\":\"${description}\",\"ratings\":${ratings}}" >/dev/null

  echo "Created: ${name} | ${genre} | ${ratings}"
}

# 1) ACTION (10)
post_movie "Action Pulse 01" "Action" "An elite courier races to prevent a citywide blackout." 7.1
post_movie "Action Pulse 02" "Action" "A retired operative returns for one final extraction mission." 7.4
post_movie "Action Pulse 03" "Action" "Street racers uncover a smuggling ring at midnight docks." 6.9
post_movie "Action Pulse 04" "Action" "A mountain rescue team battles a violent avalanche storm." 7.3
post_movie "Action Pulse 05" "Action" "A cyber unit tracks a rogue drone swarm over downtown." 7.8
post_movie "Action Pulse 06" "Action" "Two rivals join forces against a mercenary convoy." 7.0
post_movie "Action Pulse 07" "Action" "A prison break collides with a political conspiracy." 7.6
post_movie "Action Pulse 08" "Action" "A tactical squad protects a witness across hostile territory." 7.2
post_movie "Action Pulse 09" "Action" "An armored train heist turns into a rescue operation." 6.8
post_movie "Action Pulse 10" "Action" "A marine pilot fights pirates in a narrow strait." 7.5

# 2) DRAMA (10)
post_movie "Quiet Roads 01" "Drama" "Three siblings reunite to save their family bookstore." 8.1
post_movie "Quiet Roads 02" "Drama" "A teacher mentors a class through a difficult school year." 7.9
post_movie "Quiet Roads 03" "Drama" "A violinist returns home after years abroad." 7.7
post_movie "Quiet Roads 04" "Drama" "A small-town mayor faces a moral crossroads." 8.0
post_movie "Quiet Roads 05" "Drama" "A nurse balances night shifts and single parenthood." 7.8
post_movie "Quiet Roads 06" "Drama" "Old letters reveal a hidden chapter of a marriage." 8.2
post_movie "Quiet Roads 07" "Drama" "An athlete rebuilds life after career-ending injury." 7.6
post_movie "Quiet Roads 08" "Drama" "A chef reopens a closed family restaurant." 7.5
post_movie "Quiet Roads 09" "Drama" "A long train ride forces strangers to confront regrets." 7.9
post_movie "Quiet Roads 10" "Drama" "A documentary crew changes a grieving town." 8.3

# 3) COMEDY (10)
post_movie "Laugh Circuit 01" "Comedy" "A failed magician becomes a corporate team coach." 6.9
post_movie "Laugh Circuit 02" "Comedy" "Roommates fake expertise to win a cooking contest." 7.1
post_movie "Laugh Circuit 03" "Comedy" "A wedding planner juggles three ceremonies in one day." 7.0
post_movie "Laugh Circuit 04" "Comedy" "An office prank war spirals into citywide chaos." 6.8
post_movie "Laugh Circuit 05" "Comedy" "A stand-up rookie goes viral for the wrong reasons." 7.2
post_movie "Laugh Circuit 06" "Comedy" "A dog walker accidentally becomes a social media star." 6.7
post_movie "Laugh Circuit 07" "Comedy" "Two cousins run a failing motel with wild ideas." 7.3
post_movie "Laugh Circuit 08" "Comedy" "A strict librarian joins an improv club by mistake." 7.4
post_movie "Laugh Circuit 09" "Comedy" "Neighbors compete in the most dramatic garden contest." 6.9
post_movie "Laugh Circuit 10" "Comedy" "A typo in a text message changes an entire weekend." 7.0

# 4) THRILLER (10)
post_movie "Night Signal 01" "Thriller" "A radio host receives clues to a live abduction." 8.0
post_movie "Night Signal 02" "Thriller" "A forensic analyst discovers tampered evidence at home." 7.8
post_movie "Night Signal 03" "Thriller" "A ferry crossing traps passengers with a hidden attacker." 7.7
post_movie "Night Signal 04" "Thriller" "A journalist is followed after exposing a shell company." 8.1
post_movie "Night Signal 05" "Thriller" "A blackout leaves a hospital under silent siege." 7.9
post_movie "Night Signal 06" "Thriller" "A chess prodigy decodes patterns in unsolved crimes." 8.2
post_movie "Night Signal 07" "Thriller" "A mountain town hides a missing-person network." 7.6
post_movie "Night Signal 08" "Thriller" "A courier finds a drive tied to political extortion." 7.8
post_movie "Night Signal 09" "Thriller" "A late-night bus route loops through impossible streets." 7.5
post_movie "Night Signal 10" "Thriller" "A locksmith is forced to open doors for a masked crew." 7.9

# 5) HORROR (10)
post_movie "Dark Hollow 01" "Horror" "A family moves into a farmhouse with a sealed basement." 7.2
post_movie "Dark Hollow 02" "Horror" "Camp counselors hear voices from a dry lake bed." 6.9
post_movie "Dark Hollow 03" "Horror" "An antique mirror traps memories of former owners." 7.1
post_movie "Dark Hollow 04" "Horror" "A night guard patrols a museum after closing time." 7.0
post_movie "Dark Hollow 05" "Horror" "A podcast crew records sounds from an abandoned tunnel." 7.3
post_movie "Dark Hollow 06" "Horror" "A village harvest festival awakens an old curse." 6.8
post_movie "Dark Hollow 07" "Horror" "A storm reveals graves beneath a suburban street." 7.4
post_movie "Dark Hollow 08" "Horror" "A lighthouse keeper receives messages from fog." 7.2
post_movie "Dark Hollow 09" "Horror" "A board game predicts each player's final hour." 7.0
post_movie "Dark Hollow 10" "Horror" "A film archive stores footage no one remembers making." 7.5

# 6) SCI-FI (10)
post_movie "Orbit Frame 01" "Sci-fi" "A moon miner uncovers an ancient navigation map." 8.1
post_movie "Orbit Frame 02" "Sci-fi" "A time-delay signal warns Earth of a future collapse." 8.0
post_movie "Orbit Frame 03" "Sci-fi" "An android requests asylum in a floating city." 7.9
post_movie "Orbit Frame 04" "Sci-fi" "A biotech startup accidentally grows sentient coral." 7.7
post_movie "Orbit Frame 05" "Sci-fi" "Colonists on Mars lose contact after a dust wall." 8.2
post_movie "Orbit Frame 06" "Sci-fi" "A pilot wakes from stasis with missing decades." 7.8
post_movie "Orbit Frame 07" "Sci-fi" "A neural game rewrites players' memories." 7.6
post_movie "Orbit Frame 08" "Sci-fi" "A city runs on weather engines that begin to fail." 8.0
post_movie "Orbit Frame 09" "Sci-fi" "A quantum courier delivers messages between timelines." 7.9
post_movie "Orbit Frame 10" "Sci-fi" "An orbiting station shelters the last seed bank." 8.3

# 7) ROMANCE (10)
post_movie "Heartline 01" "Romance" "Two architects compete for a project then fall in love." 7.4
post_movie "Heartline 02" "Romance" "A baker and a florist share a storefront by chance." 7.1
post_movie "Heartline 03" "Romance" "A travel writer reconnects with a childhood friend." 7.5
post_movie "Heartline 04" "Romance" "A violin duet reunites former partners on stage." 7.3
post_movie "Heartline 05" "Romance" "A rainy bus stop leads to letters across cities." 7.2
post_movie "Heartline 06" "Romance" "A chef and a critic clash through anonymous reviews." 7.6
post_movie "Heartline 07" "Romance" "A dance teacher helps a pilot learn for a wedding." 7.0
post_movie "Heartline 08" "Romance" "A radio show connects callers who never met." 7.4
post_movie "Heartline 09" "Romance" "A bookstore cat keeps bringing two strangers together." 6.9
post_movie "Heartline 10" "Romance" "A missed train becomes a weekend adventure." 7.3

# 8) ANIMATION (10)
post_movie "Sketch World 01" "Animation" "A paper city comes alive when a child draws the sky." 8.0
post_movie "Sketch World 02" "Animation" "Forest creatures build a bridge before winter." 7.8
post_movie "Sketch World 03" "Animation" "A toy robot searches for its first owner." 7.9
post_movie "Sketch World 04" "Animation" "Cloud painters race to color a sunrise festival." 7.7
post_movie "Sketch World 05" "Animation" "A tiny dragon learns to breathe frost instead of fire." 7.6
post_movie "Sketch World 06" "Animation" "A shy whale writes songs in ocean currents." 8.1
post_movie "Sketch World 07" "Animation" "An inventor squirrel launches a flying acorn." 7.5
post_movie "Sketch World 08" "Animation" "A lantern spirit guides travelers at night." 7.8
post_movie "Sketch World 09" "Animation" "A paintbrush hero repairs fading memories." 8.2
post_movie "Sketch World 10" "Animation" "A village clock tower hides a portal to spring." 7.9

# 9) DOCUMENTARY (10)
post_movie "Truth Lens 01" "Documentary" "Local farmers adapt to severe water scarcity." 8.4
post_movie "Truth Lens 02" "Documentary" "An oral history of a fading railway line." 8.1
post_movie "Truth Lens 03" "Documentary" "Deep-sea divers map lost shipwreck routes." 8.0
post_movie "Truth Lens 04" "Documentary" "A profile of women rebuilding city playgrounds." 8.2
post_movie "Truth Lens 05" "Documentary" "Street musicians preserve forgotten folk songs." 7.9
post_movie "Truth Lens 06" "Documentary" "The hidden work behind emergency call centers." 8.3
post_movie "Truth Lens 07" "Documentary" "Mountain villages restore stone irrigation channels." 8.1
post_movie "Truth Lens 08" "Documentary" "Night-shift workers who keep ports running." 7.8
post_movie "Truth Lens 09" "Documentary" "Community kitchens serving disaster-hit neighborhoods." 8.5
post_movie "Truth Lens 10" "Documentary" "Young coders building low-cost assistive tools." 8.2

# 10) FANTASY (10)
post_movie "Rune Gate 01" "Fantasy" "A cartographer discovers moving borders on old maps." 8.0
post_movie "Rune Gate 02" "Fantasy" "A blacksmith forges a key that opens seasons." 7.9
post_movie "Rune Gate 03" "Fantasy" "A librarian guards books that rewrite destiny." 8.1
post_movie "Rune Gate 04" "Fantasy" "A desert caravan follows a singing compass." 7.8
post_movie "Rune Gate 05" "Fantasy" "A river spirit bargains with a stubborn prince." 7.7
post_movie "Rune Gate 06" "Fantasy" "A moon festival hides a doorway to giant gardens." 8.2
post_movie "Rune Gate 07" "Fantasy" "A masked ranger protects a city of glass." 7.9
post_movie "Rune Gate 08" "Fantasy" "A baker learns recipes that summon weather." 7.6
post_movie "Rune Gate 09" "Fantasy" "A skyship crew seeks the last thunder orchard." 8.3
post_movie "Rune Gate 10" "Fantasy" "A clockmaker pauses time for one impossible rescue." 8.1

echo "Done: 100 create-movie POST requests sent in 10 genre categories."
