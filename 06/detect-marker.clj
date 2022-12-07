;; utility functions

(defn byte-seq [^java.io.BufferedReader rdr]
  (lazy-seq
    (let [ch (.read rdr)]
      (if (= ch -1)
        '()
        (cons ch (byte-seq rdr))))))

(defn char-seq [^java.io.BufferedReader rdr]
  (map char (byte-seq rdr)))

(defn find-index [pred coll]
  (first (keep-indexed #(when (pred %2) %1) coll)))

(defn uniq? [coll]
  (= (count coll) (count (set coll))))

;; main

(def marker-length (Integer/parseInt (first *command-line-args*)))

(let [chars (char-seq (java.io.BufferedReader. *in*))]
  (println
    (+ marker-length (find-index uniq? (partition marker-length 1 chars)))))
